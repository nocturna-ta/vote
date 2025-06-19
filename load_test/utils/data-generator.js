import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export function generateVoteData() {
    return {
        voter_id : uuidv4(),
        election_pair_id: getRandomElectionPairId() ,
        region : getRandomRegion(),
        signed_transaction: generateSignedTransaction(),
    }
}

export function getRandomElectionPairId() {
    const ids = [
        'c976902f-bb0b-4103-84f5-edd74e6e928f',
        'a1234567-bb0b-4103-84f5-edd74e6e1234',
        'b2345678-bb0b-4103-84f5-edd74e6e2345'
    ]

    return ids[Math.floor(Math.random() * ids.length)];
}

export function getRandomRegion() {
    const regions = [
        'Jakarta', 'Bandung', 'Surabaya', 'Medan', 'Semarang',
        'Makassar', 'Palembang', 'Tangerang', 'Depok', 'Bekasi'
    ];
    return regions[Math.floor(Math.random() * regions.length)];
}

export function generateSignedTransaction(){
    const chars = '0123456789abcdef';
    let result = '0x';
    for (let i = 0; i < 130; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}

export function generateHeaders(){
    return {
        'Content-Type': 'application/json',
        'X-User-Id': uuidv4(),
        'X-Address-Id' : '0x' + Math.random().toString(16).substring(2, 40),
        'X-Role' : 'voter',
    };
}

export const voteIds = []