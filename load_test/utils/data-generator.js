import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

function generateOTP() {
    return Math.floor(100000 + Math.random() * 900000).toString();
}

export function generateVoteData() {
    return {
        voter_id: uuidv4(),
        election_pair_id: getRandomElectionPairId(),
        region: getRandomRegion(),
        signed_transaction: generateSignedTransaction(),
        otp_token: generateOTP(),
    };
}

// Existing functions remain unchanged
export function getRandomElectionPairId() {
    const ids = [
        'c3834ab2-7735-44b3-a4bc-7509c6c37d17',
        'd17af195-440c-4132-8169-98472fa48bb4',
        'e031d6aa-4a8d-40df-bcd8-764d9a4ba5f7'
    ];
    return ids[Math.floor(Math.random() * ids.length)];
}

export function getRandomRegion() {
    const regions = [
        'Jakarta', 'Bandung', 'Surabaya', 'Medan', 'Semarang',
        'Makassar', 'Palembang', 'Tangerang', 'Depok', 'Bekasi'
    ];
    return regions[Math.floor(Math.random() * regions.length)];
}

export function generateSignedTransaction() {
    const chars = '0123456789abcdef';
    let result = '0x';
    for (let i = 0; i < 130; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}

export function generateHeaders() {
    return {
        'Content-Type': 'application/json',
        'X-User-Id': uuidv4(),
        'X-Address-Id': '0x' + Math.random().toString(16).substring(2, 40),
        'X-Role': 'voter',
    };
}

export const voteIds = [];