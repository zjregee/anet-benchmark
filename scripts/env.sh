#!/usr/bin/env bash

scpu_cmd="taskset -c 0-3"
ccpu_cmd="taskset -c 4-11"

function benchmark() {
    for b in ${body[@]}; do
        for c in ${concurrent[@]}; do
            for ((i = 0; i < ${#repos[@]}; i++)); do
                repo=${repos[i]}
                addr="localhost:${ports[i]}" 

                svr=${repo}_server
                echo "server $svr running with $scpu_cmd"
                nohup $scpu_cmd ./output/bin/${svr} -addr="$addr" >> output/log/nohup.log 2>&1 &
                sleep 1

                cli=net_client
                echo "client $cli running with $ccpu_cmd"
                $ccpu_cmd ./output/bin/${cli} -name="$repo" -addr="$addr" -b=$b -c=$c -n=$n

                pid=$(ps -er | grep $svr | grep -v grep | awk '{print $2}')
                if [[ -n "$pid" ]]; then
                    disown $pid
                    kill -9 $pid
                fi

                sleep 1
            done
        done
    done
}
