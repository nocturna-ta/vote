import {check, sleep} from 'k6';
import {Rate} from 'k6/metrics';
import {Trend} from "k6/metrics";


export const errorRate = new Rate('errors');
export const voteCastingTrend = new Trend('vote_casting');


export function checkResponse(response, expectedStatus, action) {
    if (!response) {
        console.log(`ERROR: No response received for ${action}.`);
        return false;
    }

    if (response.status === expectedStatus) {
        return true;
    } else if (response.status === 404 && action === 'vote status') {
        console.log(`WARN: Vote status not found. Expected 200 or 404, got ${response.status}`);
        return true
    } else {
        logResponse(response, `unexpected status for ${action} expected ${expectedStatus}`)
        return false;
    }
}

export function randomSleep(min = 1, max = 3) {
    sleep(Math.random() * (max - min) + min);
}

export function logResponse(response, message) {
    if (response) {
        console.log(`ERROR: ${message} | status: ${response.status} | body: ${response.body}`);
    } else {
        console.log(`ERROR: ${message} | Response object was not provided.`);
    }
}