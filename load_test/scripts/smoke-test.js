import {config} from "../utils/config.js";
import {healthCheck} from "../scenarios/health-check.js";
import {randomSleep} from "../utils/helpers.js";
import {checkVoteStatus} from "../scenarios/vote-status.js";
import {
    completeOTPVoteFlow,
    testInvalidOTPScenarios,
    generateOTP,
    getOTPStatus
} from "../scenarios/otp-vote-flow.js";
import {generateHeaders} from "../utils/data-generator.js";

export let options = {
    stages: config.stages.smoke,
    thresholds: config.threshold,
}

export default function(){
    const headers = generateHeaders();

    // Health check
    healthCheck();
    randomSleep(0.5, 1);

    // Test complete OTP + Vote flow
    console.log('Testing complete OTP + Vote flow...');
    const otpVoteResult = completeOTPVoteFlow();

    if (otpVoteResult.success) {
        console.log(`Successfully completed OTP vote flow. Vote ID: ${otpVoteResult.voteId}`);

        // Check vote status
        randomSleep(1, 2);
        if (otpVoteResult.voteId) {
            checkVoteStatus(otpVoteResult.voteId);
        }

        // Test OTP status endpoint
        randomSleep(0.5, 1);
        getOTPStatus(otpVoteResult.voterID, headers);

    } else {
        console.log('OTP vote flow failed');
    }

    randomSleep(1, 2);

    // Test invalid OTP scenarios (should fail gracefully)
    console.log('Testing invalid OTP scenarios...');
    const voterID = generateHeaders()['X-User-Id'];

    // Generate OTP first
    const otpGenResult = generateOTP(voterID, headers);
    if (otpGenResult.success) {
        randomSleep(0.5, 1);

        // Test various invalid scenarios
        const invalidTests = testInvalidOTPScenarios(voterID, headers);

        console.log('Invalid OTP test results:', {
            wrongOTP: invalidTests.wrongOTP,
            voteWithoutOTP: invalidTests.voteWithoutOTP,
            voteWithInvalidOTP: invalidTests.voteWithInvalidOTP
        });
    }

    randomSleep(1, 2);
}