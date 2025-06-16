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
        http_reqs: ['rate>50'],
    }
}

export default function (){
    const userJourney = Math.random();

    if(userJourney < 0.1){
        healthCheck()
    }else if (userJourney < 0.6) {
        const {voteId} = castVote()
        randomSleep(2, 5)

        if(voteId){
            checkVoteStatus(voteId);
            randomSleep(1, 3);

            if (Math.random() < 0.3){
                checkVoteStatus(voteId)
            }
        }
    }else{
        checkVoteStatus()
    }

    randomSleep(2,8);
}