import {generateHeaders, generateVoteData} from "../utils/data-generator.js";
import http from 'k6/http';
import {config} from "../utils/config.js";
import {checkResponse, logResponse} from "../utils/helpers.js";

export function castVote(){
    const voteData = generateVoteData()
    const headers = generateHeaders()

    const response = http.post(
        `${config.base_url}/v1/vote/cast`,
        JSON.stringify(voteData),
        { headers }
    );

    const success = checkResponse(response, 200, 'cast vote');

    if (!success){
        logResponse(response, 'Cast Vote Failed')
    }

    let voteId = null
    if(response.status === 200) {
        try {
            const responseBody = JSON.parse(response.body);
            voteId = responseBody.data?.id
        }catch (e){
            console.log('Failed to parse cast vote response')
        }
    }

    return {response, voteId}
}