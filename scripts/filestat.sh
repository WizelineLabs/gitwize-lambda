# #!/usr/bin/env bash
# TODO support time range
stat() {
    vRepo=$1
    vRepoPath=$2
    cd $vRepoPath
    export PATH=$PATH:"$(pwd)/scripts"
    git log --pretty=format:"%H|%ae|%an|%ad" --date=unix --no-merges | \
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
        values=$(IFS="," ; echo "${value_arr[*]}")
        sql="INSERT INTO $GW_DB_NAME.file_stat_data (repository_id, hash, author_email, author_name,
                file_name, addition_loc, deletion_loc, year, month, day, hour, commit_time_stamp, modification_loc)
            VALUES $values
            ON DUPLICATE KEY UPDATE addition_loc=VALUES(addition_loc), deletion_loc=VALUES(deletion_loc), modification_loc=VALUES(modification_loc)"

        # Note for local: make sure mysql is on $PATH, aliases not working in non-interactive shells
        mysql -h$GW_DB_HOST -P$GW_DB_PORT -u$GW_DB_USER -p$GW_DB_PASSWORD -e "$sql"
    )
}

stat $1 $2