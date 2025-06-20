import { config } from '../utils/config.js';
import { healthCheck } from '../scenarios/health-check.js';
import { randomSleep } from '../utils/helpers.js';
import {castVote} from "../scenarios/vote-casting.js";

export let options = {
    stages: config.stages.load,
    thresholds: {
        ...config.thresholds,
        http_req_duration: ['p(95)<1000'], // Relaxed threshold for load test
    }
};

export default function() {
    const scenario = Math.random();

        if (scenario < 0.8){
            const {voteId} = castVote()
            randomSleep(1, 3)
        }else{
            healthCheck();
        }


    randomSleep(1, 3);
}