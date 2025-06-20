import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import {generateHeaders} from "../utils/data-generator.js";
import http from 'k6/http';
import {checkResponse, logResponse} from "../utils/helpers.js";
import {config} from "../utils/config.js";

export function checkVoteStatus(voteId = null){
    const id = voteId || uuidv4();
    const headers = generateHeaders()

    const response = http.get(
        `${config.base_url}/v1/vote/${id}/status`,
        { headers }
    );

    const expectedStatus = voteId ? 200 : 400;
    const success = checkResponse(response, expectedStatus, 'vote status check');

    if(!success){
        logResponse(response, 'Vote Status Check Failed');
    }

    return response
}