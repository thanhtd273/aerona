package com.thanhtd.aerona.payment.service.impl;

import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;
import com.thanhtd.aerona.base.service.RedisService;
import com.thanhtd.aerona.payment.dto.PaymentInfo;
import com.thanhtd.aerona.payment.dto.PaymentSource;
import com.thanhtd.aerona.payment.dto.PaypalTokenResponse;
import com.thanhtd.aerona.payment.dto.paypal.OrderResponse;
import com.thanhtd.aerona.payment.service.PaypalService;
import lombok.RequiredArgsConstructor;

import org.apache.commons.lang3.ObjectUtils;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.util.MultiValueMap;
import org.springframework.web.client.RestTemplate;

import java.util.Base64;

@Service
@RequiredArgsConstructor
public class PaypalServiceImpl implements PaypalService {

    @Value("${paypal.client-id}")
    private String clientId;

    @Value("${paypal.client-secret}")
    private String clientSecret;

    @Value("${paypal.base.url}")
    private String baseUrl;

    private final RestTemplate restTemplate;

    private final RedisService redisService;

    @Override
    public String getAccessToken() throws LogicException {
        String redisKey = String.format("paypal_access_token#%s", clientId);
        String accessToken = redisService.getValue(redisKey);
        if (!ObjectUtils.isEmpty(accessToken)) {
            return accessToken;
        }
        String auth = String.format("%s:%s", clientId, clientSecret);
        String encodeAuth = Base64.getEncoder().encodeToString(auth.getBytes());

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_FORM_URLENCODED);
        headers.set("Authorization", "Basic " + encodeAuth);

        MultiValueMap<String, String> body = new LinkedMultiValueMap<>();
        body.add("grant_type", "client_credentials");

        HttpEntity<MultiValueMap<String, String>> request = new HttpEntity<>(body, headers);
        ResponseEntity<PaypalTokenResponse> response = restTemplate.postForEntity(String.format("%s/v1/oauth2/token", baseUrl), request, PaypalTokenResponse.class);
        PaypalTokenResponse data = response.getBody();
        if (ObjectUtils.isEmpty(data)) {
            throw new LogicException(ErrorCode.FAIL);
        }
        accessToken = data.getAccessToken();
        redisService.save(redisKey, accessToken, data.getExpiresIn());
        return accessToken;
    }

    @Override
    public OrderResponse createOrder(PaymentInfo paymentInfo) throws LogicException {
        String accessToken = getAccessToken();
        String createOrderUrl = String.format("%s/v2/checkout/orders", baseUrl);
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.setBearerAuth(accessToken);
        HttpEntity<String> request = getStringHttpEntity(String.format("""
                {
                    "intent": "CAPTURE",
                    "purchase_units": [
                        {
                            "amount": {
                                "currency_code": "%s",
                                "value": "%.2f"
                            }
                        }
                    ]
                }
                """, paymentInfo.getCurrency(), paymentInfo.getAmount()), headers);
        ResponseEntity<OrderResponse> response = restTemplate.postForEntity(createOrderUrl, request, OrderResponse.class);

        return response.getBody();
    }

    private static HttpEntity<String> getStringHttpEntity(String paymentInfo, HttpHeaders headers) {
        String orderPayload = paymentInfo;
        HttpEntity<String> request = new HttpEntity<>(orderPayload, headers);
        return request;
    }

    @Override
    public OrderResponse captureOrder(String orderId) throws LogicException {
        String accessToken = getAccessToken();
        String captureUrl = String.format("%s/v2/checkout/orders/%s/capture", baseUrl, orderId);
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.setBearerAuth(accessToken);

        HttpEntity<String> request = new HttpEntity<>("", headers);
        ResponseEntity<OrderResponse> response = restTemplate.exchange(captureUrl, HttpMethod.POST, request, OrderResponse.class);
        return response.getBody();
    }

    @Override
    public OrderResponse confirmPaymentSource(String orderId, PaymentSource paymentSource)
            throws LogicException {
        String accessToken = getAccessToken();
        String confirmPaymentUrl = String.format("%s/v2/checkout/orders/%s/confirm-payment-source", baseUrl, orderId);
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.setBearerAuth(accessToken);
        String payload = String.format("""
                {
                    "payment_source": {
                        "card": {
                            "number": "%s",
                            "expiry": "%s"
                        }
                    }
                }
                """, paymentSource.getNumber(), paymentSource.getExpiry());
        HttpEntity<String> request = new HttpEntity<>(payload, headers);
        ResponseEntity<OrderResponse> response = restTemplate.postForEntity(confirmPaymentUrl, request, OrderResponse.class);

        return response.getBody();
    }
}
