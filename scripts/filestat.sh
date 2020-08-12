#!/usr/bin/env bash

# This script looks at modified files in a commit and tries to determine below metrics for every file in every commit.
# `Addition`
# `Deletion`
# `Modification`
# `Churn` - if the line was added/modified less than 21 days from current commit
# `Refactoring` - otherwise
#
# Run:
#   ./filestat.sh <repo_id> <repo_path>
# Output:
#     Data will be persisted to `file_stat_data` table according to setting of below vars:
#    - GW_DB_HOST
#    - GW_DB_PORT
#    - GW_DB_USER
#    - GW_DB_PASSWORD
#    - GW_DB_NAME
#
#    Enable `GW_DEBUG` for more info.

function log_debug() {
    if [ "$GW_DEBUG" = "true" ]; then
        echo "[DEBUG] $*"
    fi
}

# Customized from https://stackoverflow.com/a/32616440
# Massage the @@ counts so they are usable
function decomma_counts() {
    cat | awk -F',' 'BEGIN { convert = 0; }
                    /^@@ / { convert=1; }
                    /^/  {
                            if ( convert == 1 ) { print $1,$2,$3;
                            } else { print $0; }
                    convert=0;
                }'
}

# Extract all new changes added with the line count
function numerize_lines() {
    cat | awk 'BEGIN { display=0; line=0; left=0; out=1;}
                        /^@@ / { out=0; inc=0; line=-$2; line--; display=line; left=line;        }
                        /^[-]/   { left++; display=left; inc=0; }
                        /^[+]/   { line++; display=line; inc=0; }
                        /^[-+][-+][-+] / { out=0; inc=0; }
                        /^/    { 
                                    line += inc;
                                    left += inc;
                                    display += inc;
                                    if ( out == 1 ) {
                                        print display,$0;
                                    } else {
                                        print $0;
                                    }
                                    out = 1;
                                    inc = 1;
                                    display = line;
                                }'
}

# Shows last updated data a line in a file before a commit
# `blame <line_no> <commit_hash> <file>`
function blame() {
    vLine=$1
    vCommit="$2^"
    vFile=$3
    git blame --line-porcelain -L $vLine,+1 "$vCommit" --date=unix -- "$vFile"
}

# Gets last editting timestamp of a line before a commit
# `last_editted_time_before_commit <line_no> <commit_hash> <file>`
function last_editted_time_before_commit() {
    blame $1 $2 $3 | grep author-time | cut -d" " -f 2
}

# Lists all changed file in a commit
# `list_file_changed <commit_hash>`
function list_file_changed() {
    vCommit=$1
    git show $vCommit --numstat --pretty=oneline | sed -n '2,$'p | awk '{ print $3 }' |
        awk '!/.*.json/ { print }' # ignore *lock.json
}

# Calculates metrics for a single file in a commit
# `filestat_single_file <commit_hash> <file>`
function filestat_single_file() {
    vCommit=$1
    vFile=$2

    # additions, deletions
    vNumstat=$(git show --numstat $vCommit -- $vFile | sed -n '$'p)
    vAdditions=$(echo "$vNumstat" | awk '{ print $1 }')
    vDeletions=$(echo "$vNumstat" | awk '{ print $2 }')

    # modifications
    # count for lines which contains both additions ({+...+}) and deletions ([-...-]) provided by `word-diff` option
    vModifications=$(git show --word-diff $vCommit -- $vFile | grep "{+.*+}" | grep "\[-.*-\]" | wc -l)
    vNewCode=$(($vAdditions - $vModifications))
    vDeletions=$(($vDeletions - $vModifications))

    # churn/refactor
    vChanges=$(git show -U0 $vCommit -- $vFile | grep -v "^+" | decomma_counts | numerize_lines)

    vChurn=0
    vRefactor=0
    echo "$vChanges" | awk '/^[0-9]+ -.*/ { print }' |
        cut -d" " -f 1 |
        (
            while read vLineNo; do
                vCommitTs=$(git log --pretty=format:'%H|%ad' --date=unix | grep $vCommit | cut -d"|" -f 2)
                vLastEdittedTs=$(last_editted_time_before_commit $vLineNo $vCommit $vFile)
                days=$((($vCommitTs - $vLastEdittedTs) / 24 / 3600))

                # If last change is not old enough (21 days), this change of the line is considered `churn`
                if [ "$days" -lt 22 ]; then
                    vChurn=$(($vChurn + 1))
                else
                    vRefactor=$(($vRefactor + 1))
                fi
            done
            echo "$vCommit,$vFile,$vNewCode,$vDeletions,$vModifications,$vChurn,$vRefactor"
        )
}

# Calculates metrics for every file in a commit
# `filestat_single_commit <commit_hash>`
function filestat_single_commit() {
    vCommit=$1
    list_file_changed $vCommit |
        while read L; do
            filestat_single_file $vCommit $L
        done
}

# Calculates metrics for every file in every commit of a repo
# `filestat <repo_ID> <repo_path>`
function filestat() {
    vRepo=$1
    vRepoPath=$2
    vFromDate=$3
    vToDate=$4

    vRange=""
    if [[ "$vFromDate" ]]; then
        vRange="$vRange --since '$vFromDate'"
    fi
    if [[ "$vToDate" ]]; then
        vRange="$vRange --until '$vToDate'"
    fi

    export PATH=$PATH:"$(pwd)/scripts"

    echo "Updating file_stat_data for repo=$vRepo within range [$vFromDate - $vToDate]"
    cd $vRepoPath

    commit_cnt=0
    eval "git log --pretty=format:'%H|%ae|%an|%ad' --date=unix --no-merges $vRange" |
        (
            value_arr=()
            while IFS="|" read vHash vEmail vName vTimestamp; do
                commit_cnt=$(($commit_cnt+1))
                vCommit=$vHash
                log_debug "Running stats for commit: filestat_single_commit $vCommit"
                filestats=$(filestat_single_commit $vCommit)
                log_debug "----- Result -----"
                log_debug "\n$filestats"
                log_debug "----- End Result -----"

                vParsedDate="FROM_UNIXTIME($vTimestamp)"
                while IFS="," read vHash vFilename vAdditions vDeletions vModifications vChurn vRefactoring; do
                    value_arr+=("($vRepo,'$vHash','$vEmail','$vName','$vFilename',$vAdditions,$vDeletions,$vModifications,$vChurn,$vRefactoring,YEAR($vParsedDate),MONTH($vParsedDate),DAY($vParsedDate),HOUR($vParsedDate),$vParsedDate)")
                done <<<"$filestats"

                if [ "$commit_cnt" -eq 10 ]; then
                    insert_data
                    commit_cnt=0
                fi
                
            done
        )
}

function insert_data() {
    # join SQL insert values
    values=$(
        IFS=","
        echo "${value_arr[*]}"
    )

    if [[ -z "$values" ]]; then
        echo "No commits found. Do nothing!"
        exit 0
    fi

    sql="INSERT INTO $GW_DB_NAME.file_stat_data (repository_id, hash, author_email, author_name,
                file_name, addition_loc, deletion_loc, modification_loc, churn_cnt, refactoring_cnt, year, month, day, hour, commit_time_stamp)
            VALUES $values
            ON DUPLICATE KEY UPDATE addition_loc=VALUES(addition_loc),
                deletion_loc=VALUES(deletion_loc),
                modification_loc=VALUES(modification_loc),
                churn_cnt=VALUES(churn_cnt),
                refactoring_cnt=VALUES(refactoring_cnt)"

    echo "Executing SQL: $sql"
    # # Note for local: make sure mysql is on $PATH, aliases not working in non-interactive shells
    mysql -h$GW_DB_HOST -P$GW_DB_PORT -u$GW_DB_USER -p$GW_DB_PASSWORD -e "$sql"
}

filestat $1 $2 $3 $4
