import {config} from "../utils/config.js";
import {castVote} from "../scenarios/vote-casting.js";
import {checkVoteStatus} from "../scenarios/vote-status.js";
import {randomSleep} from "../utils/helpers.js";

export let options = {
    stages: config.stages.spike,
    thresholds: {
        http_req_duration: ['p(95)<3000'],
        http_req_failed: ['rate<0.3'],
    }
}

export default function (){
    if (Math.random() < 0.8){
        castVote()
    }else{
        checkVoteStatus()
    }

    randomSleep(0.1, 1)
}