package com.thanhtd.aerona.payment.config;

import lombok.extern.slf4j.Slf4j;
import org.apache.hc.client5.http.ConnectTimeoutException;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.client.ClientHttpResponse;
import org.springframework.util.StreamUtils;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.ResponseErrorHandler;

import java.io.IOException;
import java.net.URI;
import java.nio.charset.StandardCharsets;

@Slf4j
public class HttpResponseErrorHandler implements ResponseErrorHandler {

    @Override
    public boolean hasError(ClientHttpResponse response) throws IOException {
        try {
            HttpStatusCode statusCode = response.getStatusCode();
            return (statusCode.value() == 606 ||
                    statusCode.is4xxClientError() ||
                    statusCode.is5xxServerError());
        } catch (IOException e) {
            log.error("Failed to read response status", e);
            return true;
        }
    }

    @Override
    public void handleError(URI url, HttpMethod method, ClientHttpResponse response) throws IOException {
        HttpStatusCode statusCode;
        String statusText;
        String responseBody;

        try {
            statusCode = response.getStatusCode();
            statusText = response.getStatusText();
            responseBody = StreamUtils.copyToString(response.getBody(), StandardCharsets.UTF_8);
        } catch (IOException e) {
            log.error("Failed to process error response for {} {}", method, url, e);
            throw e;
        }

        log.error("API call failed for {} {} - Status: {} - {}, Response: {}",
                method, url, statusCode.value(), statusText, responseBody);

        if (statusCode.value() == 606) {
            throw new ConnectTimeoutException("API timeout (606) for " + url);
        } else if (statusCode.value() == 429) {
            log.warn("Rate limit exceeded for {} {}", method, url);
            throw new HttpClientErrorException(statusCode, statusText, responseBody.getBytes(), StandardCharsets.UTF_8);
        } else if (statusCode.is4xxClientError()) {
            throw new HttpClientErrorException(statusCode, statusText, responseBody.getBytes(), StandardCharsets.UTF_8);
        } else if (statusCode.is5xxServerError()) {
            throw new HttpServerErrorException(statusCode, statusText, responseBody.getBytes(), StandardCharsets.UTF_8);
        }
    }
}