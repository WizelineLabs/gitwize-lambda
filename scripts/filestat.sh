# #!/usr/bin/env bash
stat() {
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
    eval "git log --pretty=format:'%H|%ae|%an|%ad' --date=unix --no-merges $vRange" | \
    (
        value_arr=()
        while IFS="|" read vHash vEmail vName vTimestamp
        do
            stats=$(git show $vHash | diffstat -m -t | grep '^[a-f0-9]*,[0-9]*,.*$')
            while IFS="," read vAdditions vDeletions vModifications vFilename
            do
                vFilename=$(echo $vFilename | sed 's|"||g')
                vParsedDate="FROM_UNIXTIME($vTimestamp)"
                value_arr+=("($vRepo,'$vHash','$vEmail','$vName','$vFilename',$vAdditions,$vDeletions,YEAR($vParsedDate),MONTH($vParsedDate),DAY($vParsedDate),HOUR($vParsedDate),$vParsedDate,$vModifications)")
            done <<< "$stats"
        done

        # join SQL insert values
        values=$(IFS="," ; echo "${value_arr[*]}")

        if [[ -z "$values" ]]; then
            echo "No commits found. Do nothing!"
            exit 0
        fi

        sql="INSERT INTO $GW_DB_NAME.file_stat_data (repository_id, hash, author_email, author_name,
                file_name, addition_loc, deletion_loc, year, month, day, hour, commit_time_stamp, modification_loc)
            VALUES $values
            ON DUPLICATE KEY UPDATE addition_loc=VALUES(addition_loc), deletion_loc=VALUES(deletion_loc), modification_loc=VALUES(modification_loc)"
        
        # Note for local: make sure mysql is on $PATH, aliases not working in non-interactive shells
        mysql -h$GW_DB_HOST -P$GW_DB_PORT -u$GW_DB_USER -p$GW_DB_PASSWORD -e "$sql"
    )
}

stat $1 $2 $3 $4
