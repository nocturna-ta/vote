import {castVote} from "../scenarios/vote-casting.js";
import {checkVoteStatus} from "../scenarios/vote-status.js";
import {healthCheck} from "../scenarios/health-check.js";
import {randomSleep} from "../utils/helpers.js";
import {config} from "../utils/config.js";

export let options = {
    stages: config.stages.stress,
    thresholds : {
        http_req_duration: ['p(95)<2000'],
        http_req_failed: ['rate<0.2'],
        http_reqs:['rate>200'],
    }
}

export default function (){
    const actions = Math.floor(Math.random() * 3) + 1;

    for (let i = 0; i < actions; i++){
        const scenario = Math.random()

        if (scenario < 0.7){
            castVote()
        }else if (scenario < 0.9){
            checkVoteStatus()
        }else{
            healthCheck()
        }

        randomSleep(0.1, 0.5)

    }
}