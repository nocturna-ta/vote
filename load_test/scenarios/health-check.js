import http from 'k6/http';
import { config } from '../utils/config.js';
import { checkResponse } from '../utils/helpers.js';

export function healthCheck() {
    const response = http.get(`${config.base_url}/health`);

    checkResponse(response, 200, 'health check');

    return response;
}