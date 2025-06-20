import {config} from "../utils/config.js";
import {healthCheck} from "../scenarios/health-check.js";
import {randomSleep} from "../utils/helpers.js";
import {checkVoteStatus} from "../scenarios/vote-status.js";
import {castVote} from "../scenarios/vote-casting.js";

export let options = {
    stages: config.stages.smoke,
    thresholds: config.threshold,
}

export default function(){
    healthCheck()
    randomSleep(0.5, 1);

    const {voteId} = castVote();
    randomSleep(1, 2)

    if (voteId){
        checkVoteStatus(voteId)
    }else{
        checkVoteStatus()
    }

    randomSleep(1, 2)
}