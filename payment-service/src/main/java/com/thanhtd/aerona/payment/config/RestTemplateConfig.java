package com.thanhtd.aerona.payment.config;

import lombok.extern.slf4j.Slf4j;
import org.apache.hc.client5.http.HttpRequestRetryStrategy;
import org.apache.hc.client5.http.config.ConnectionConfig;
import org.apache.hc.client5.http.config.RequestConfig;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.HttpClientBuilder;
import org.apache.hc.client5.http.impl.io.PoolingHttpClientConnectionManager;
import org.apache.hc.core5.http.HttpRequest;
import org.apache.hc.core5.http.HttpResponse;
import org.apache.hc.core5.http.io.SocketConfig;
import org.apache.hc.core5.http.protocol.HttpContext;
import org.apache.hc.core5.util.TimeValue;
import org.apache.hc.core5.util.Timeout;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.MediaType;
import org.springframework.http.client.*;
import org.springframework.http.converter.StringHttpMessageConverter;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.web.client.RestTemplate;

import java.io.IOException;
import java.net.ConnectException;
import java.net.SocketTimeoutException;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;

@Configuration
@Slf4j
public class RestTemplateConfig {

    private static final int MAX_RETRY = 3;
    private static final int TIMEOUT_SECONDS = 5;
    private static final int MAX_TOTAL_CONNECTIONS = 100;
    private static final int MAX_PER_ROUTE = 20;

    @Bean
    public RestTemplate restTemplate() {
        RestTemplate restTemplate = new RestTemplate(getClientHttpRequestFactory());
        restTemplate.getInterceptors().add(new CustomClientHttpRequestInterceptor());
        restTemplate.setErrorHandler(new HttpResponseErrorHandler());
        MappingJackson2HttpMessageConverter converter = new MappingJackson2HttpMessageConverter();
        converter.setSupportedMediaTypes(Arrays.asList(MediaType.APPLICATION_JSON, MediaType.APPLICATION_OCTET_STREAM));
        restTemplate.getMessageConverters().add(converter);
        restTemplate.getMessageConverters().forEach(conv -> {
            if (conv instanceof StringHttpMessageConverter stringConv) {
                stringConv.setDefaultCharset(StandardCharsets.UTF_8);
            }
        });
        return restTemplate;
    }

    private ClientHttpRequestFactory getClientHttpRequestFactory() {
        int timeoutMs = TIMEOUT_SECONDS * 1000;

        ConnectionConfig connectionConfig = ConnectionConfig.custom()
                .setConnectTimeout(Timeout.ofMilliseconds(timeoutMs))
                .build();
        RequestConfig requestConfig = RequestConfig.custom()
                .setConnectionRequestTimeout(Timeout.ofMilliseconds(timeoutMs))
                .build();
        SocketConfig socketConfig = SocketConfig.custom()
                .setSoTimeout(Timeout.ofMilliseconds(timeoutMs))
                .build();

        PoolingHttpClientConnectionManager connectionManager = new PoolingHttpClientConnectionManager();
        connectionManager.setMaxTotal(MAX_TOTAL_CONNECTIONS);
        connectionManager.setDefaultMaxPerRoute(MAX_PER_ROUTE);
        connectionManager.setDefaultConnectionConfig(connectionConfig);
        connectionManager.setDefaultSocketConfig(socketConfig);

        CloseableHttpClient httpClient = HttpClientBuilder.create()
                .setConnectionManager(connectionManager)
                .setDefaultRequestConfig(requestConfig)
                .evictExpiredConnections()
                .evictIdleConnections(Timeout.ofSeconds(30))
                .setRetryStrategy(new CustomHttpRequestRetryStrategy(MAX_RETRY))
                .build();

        HttpComponentsClientHttpRequestFactory factory = new HttpComponentsClientHttpRequestFactory(httpClient);
        factory.setReadTimeout(timeoutMs);
        return factory;
    }

    private record CustomHttpRequestRetryStrategy(int maxRetry) implements HttpRequestRetryStrategy {

        @Override
            public boolean retryRequest(HttpRequest request, IOException exception, int executionCount, HttpContext context) {
                if (executionCount > maxRetry) return false;
                return exception instanceof ConnectException || exception instanceof SocketTimeoutException;
            }

            @Override
            public boolean retryRequest(HttpResponse response, int executionCount, HttpContext context) {
                if (executionCount > maxRetry) return false;
                int responseCode = response.getCode();
                if (responseCode >= 500 && responseCode < 600) {
                    log.info("Retrying for 5xx code: {}", responseCode);
                    return true;
                } else if (responseCode == 429) {
                    log.info("Too many requests (429), retrying...");
                    return true;
                }
                return false;
            }

            @Override
            public TimeValue getRetryInterval(HttpResponse response, int executionCount, HttpContext context) {
                int baseDelayMs = 100;
                int delay = baseDelayMs * (1 << (executionCount - 1));
                return TimeValue.ofMilliseconds(Math.min(delay, 5000));
            }
        }

}
