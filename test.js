import http from 'k6/http';
import { check, sleep } from 'k6';

// Load the certificate and key files in the global scope
const cert = open('/Users/anle/Documents/Project/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem');
const key = open('/Users/anle/Documents/Project/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/32b13cdd7918b000ba163fa97a6f3168e3563e2b499320b9bebae699dc0d370e_sk');

// Print out the data for debugging purposes
// console.log('Certificate Data:', cert);
// console.log('Key Data:', key);

export let options = {
    stages: [
        { duration: '15s', target: 400 }, // Ramp-up to N users over X second and stay at Y users
    ],
};

export default function () {
    // Randomly select an asset from asset1 to asset6
    const assetId = Math.floor(Math.random() * 6) + 1; // This will generate a number between 1 and 6
    const url = `http://localhost:8080/assets/asset${assetId}`;

    // Create the form data
    const formData = {
        cert: http.file(cert, 'cert.pem'),
        key: http.file(key, 'key.pem'),
    };

    // Send the POST request with the form data
    const res = http.post(url, formData);

    // Check the response
    check(res, {
        'is status 200': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
    });

    sleep(1);
}