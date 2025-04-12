package com.thanhtd.aerona.booking.config;

import org.apache.http.conn.ConnectTimeoutException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.MediaType;
import org.springframework.http.client.ClientHttpResponse;
import org.springframework.http.converter.StringHttpMessageConverter;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.web.client.ResponseErrorHandler;
import org.springframework.web.client.RestTemplate;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;

@Configuration
public class RestTemplateConfig {

    private static final Logger logger = LoggerFactory.getLogger(RestTemplateConfig.class);

    @Bean
    public RestTemplate restTemplate() {
        RestTemplate restTemplate = new RestTemplate();
        restTemplate.getInterceptors().add(new CustomClientHttpRequestInterceptor());
        restTemplate.setErrorHandler(errorHandler());
        MappingJackson2HttpMessageConverter converter = new MappingJackson2HttpMessageConverter();
        converter.setSupportedMediaTypes(Arrays.asList(MediaType.APPLICATION_JSON, MediaType.APPLICATION_OCTET_STREAM));
        restTemplate.getMessageConverters().add(converter);
        restTemplate.getMessageConverters().forEach(httpMessageConverter -> {
            if(httpMessageConverter instanceof StringHttpMessageConverter stringHttpMessageConverter) {
                stringHttpMessageConverter.setDefaultCharset(StandardCharsets.UTF_8);
            }
        });
        return restTemplate;
    }

    @Bean
    public CustomClientHttpRequestInterceptor customClientHttpRequestInterceptor() {
        return new CustomClientHttpRequestInterceptor();
    }

    ResponseErrorHandler errorHandler() {
        return new ResponseErrorHandler() {
            @Override
            public boolean hasError(ClientHttpResponse response) throws IOException {

                return (response.getStatusCode().value() == 606 ||
                        response.getStatusCode().is4xxClientError() ||
                        response.getStatusCode().is5xxServerError());
            }

            @Override
            public void handleError(ClientHttpResponse response) throws IOException {
                HttpStatusCode statusCode = response.getStatusCode();
                String statusText = response.getStatusText();

                logger.error("API call failed with status: {} - {}", statusCode.value(), statusText);
                if(response.getStatusCode().value() == 606 ){
                    throw new ConnectTimeoutException();
                }
            }
        };
    }
}
