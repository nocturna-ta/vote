import {config} from "../utils/config.js";
import {healthCheck} from "../scenarios/health-check.js";
import {castVote} from "../scenarios/vote-casting.js";
import {randomSleep} from "../utils/helpers.js";
import {checkVoteStatus} from "../scenarios/vote-status.js";

export let options = {
    stages: config.stages.endurance,
    thresholds: {
        http_req_duration :['p(95)<1500'],
        http_req_failed: ['rate<0.15'],
        http_reqs: ['rate>40'],
    }
}

export default function (){
    const userJourney = Math.random();

    if(userJourney < 0.1){
        healthCheck()
    }else if (userJourney < 0.6) {
        castVote()
        randomSleep(0.1, 0.3)
    }

    randomSleep(0.05,0.3);
}