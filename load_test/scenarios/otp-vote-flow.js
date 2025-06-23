import {
    generateHeaders,
    generateOTPRequestData,
    generateOTPVerificationData,
    generateVoteData,
    otpTokens
} from "../utils/data-generator.js";
import http from 'k6/http';
import {config} from "../utils/config.js";
import {checkResponse, logResponse} from "../utils/helpers.js";

export function completeOTPVoteFlow() {
    const headers = generateHeaders();
    const voterID = generateOTPRequestData().voter_id;

    // Step 1: Generate OTP
    const otpGenResponse = generateOTP(voterID, headers);
    if (!otpGenResponse.success) {
        logResponse(otpGenResponse.response, 'OTP Generation Failed');
        return { success: false, voteId: null };
    }

    // Step 2: Verify OTP (simulate user entering correct OTP)
    const otpVerifyResponse = verifyOTP(voterID, "123456", headers); // Use test OTP
    if (!otpVerifyResponse.success || !otpVerifyResponse.token) {
        logResponse(otpVerifyResponse.response, 'OTP Verification Failed');
        return { success: false, voteId: null };
    }

    // Step 3: Cast vote with OTP token
    const voteResponse = castVoteWithOTP(voterID, otpVerifyResponse.token, headers);
    if (!voteResponse.success) {
        logResponse(voteResponse.response, 'Vote Casting Failed');
        return { success: false, voteId: null };
    }

    return {
        success: true,
        voteId: voteResponse.voteId,
        voterID: voterID
    };
}

export function generateOTP(voterID, headers) {
    const otpData = {
        voter_id: voterID,
        purpose: "vote_cast"
    };

    const response = http.post(
        `${config.base_url}/v1/otp/generate`,
        JSON.stringify(otpData),
        { headers }
    );

    const success = checkResponse(response, 200, 'generate OTP');

    return { response, success };
}

export function verifyOTP(voterID, otpCode, headers) {
    const verifyData = {
        voter_id: voterID,
        purpose: "vote_cast",
        code: otpCode
    };

    const response = http.post(
        `${config.base_url}/v1/otp/verify`,
        JSON.stringify(verifyData),
        { headers }
    );

    const success = checkResponse(response, 200, 'verify OTP');

    let token = null;
    if (success && response.status === 200) {
        try {
            const responseBody = JSON.parse(response.body);
            if (responseBody.data && responseBody.data.is_valid) {
                token = responseBody.data.otp_token;
                // Store token for potential reuse in testing
                otpTokens.set(voterID, token);
            }
        } catch (e) {
            console.log('Failed to parse OTP verification response:', e);
        }
    }

    return { response, success, token };
}

export function resendOTP(voterID, headers) {
    const resendData = {
        voter_id: voterID,
        purpose: "vote_cast"
    };

    const response = http.post(
        `${config.base_url}/v1/otp/resend`,
        JSON.stringify(resendData),
        { headers }
    );

    const success = checkResponse(response, 200, 'resend OTP');

    return { response, success };
}

export function getOTPStatus(voterID, headers) {
    const response = http.get(
        `${config.base_url}/v1/otp/status?voter_id=${voterID}&purpose=vote_cast`,
        { headers }
    );

    const success = checkResponse(response, 200, 'get OTP status');

    return { response, success };
}

export function castVoteWithOTP(voterID, otpToken, headers) {
    const voteData = generateVoteData();
    voteData.voter_id = voterID; // Use the same voter ID
    voteData.otp_token = otpToken; // Use the verified OTP token

    const response = http.post(
        `${config.base_url}/v1/vote/cast`,
        JSON.stringify(voteData),
        { headers }
    );

    const success = checkResponse(response, 200, 'cast vote with OTP');

    let voteId = null;
    if (success && response.status === 200) {
        try {
            const responseBody = JSON.parse(response.body);
            voteId = responseBody.data?.id;
        } catch (e) {
            console.log('Failed to parse vote response:', e);
        }
    }

    return { response, success, voteId };
}

// Test invalid OTP scenarios
export function testInvalidOTPScenarios(voterID, headers) {
    // Test with wrong OTP code
    const wrongOTPResponse = verifyOTP(voterID, "000000", headers);

    // Test vote without OTP token
    const voteWithoutOTPData = generateVoteData();
    voteWithoutOTPData.voter_id = voterID;
    delete voteWithoutOTPData.otp_token; // Remove OTP token

    const voteWithoutOTPResponse = http.post(
        `${config.base_url}/v1/vote/cast`,
        JSON.stringify(voteWithoutOTPData),
        { headers }
    );

    // Test vote with invalid OTP token
    const voteWithInvalidOTPData = generateVoteData();
    voteWithInvalidOTPData.voter_id = voterID;
    voteWithInvalidOTPData.otp_token = "invalid_token";

    const voteWithInvalidOTPResponse = http.post(
        `${config.base_url}/v1/vote/cast`,
        JSON.stringify(voteWithInvalidOTPData),
        { headers }
    );

    return {
        wrongOTP: checkResponse(wrongOTPResponse.response, 200, 'wrong OTP verification'),
        voteWithoutOTP: checkResponse(voteWithoutOTPResponse, 400, 'vote without OTP'),
        voteWithInvalidOTP: checkResponse(voteWithInvalidOTPResponse, 401, 'vote with invalid OTP')
    };
}