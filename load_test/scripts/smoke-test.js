import {config} from "../utils/config.js";
import {healthCheck} from "../scenarios/health-check.js";
import {randomSleep} from "../utils/helpers.js";
import {checkVoteStatus} from "../scenarios/vote-status.js";
import {castVote} from "../scenarios/vote-casting.js";

export let options = {
    stages: config.stages.smoke,
    thresholds:{
        ...config.threshold,
        'http_req_duration': ['p(95)<500'], // Relaxed for smoke
        'http_req_failed': ['rate<0.01'],
        'http_reqs': ['rate>5']
    },
}

export default function(){
    healthCheck();
    randomSleep(0.1, 0.3);

    const { voteId } = castVote();
    randomSleep(0.1, 0.3);

    if (voteId) {
        checkVoteStatus(voteId);
    } else {
        checkVoteStatus("sample-vote-id");
    }
    randomSleep(0.1, 0.3);
}