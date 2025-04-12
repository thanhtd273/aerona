package com.thanhtd.aerona.payment.service.impl;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.thanhtd.aerona.payment.constant.TopicConstant;
import com.thanhtd.aerona.payment.dto.BookingInfo;
import com.thanhtd.aerona.payment.service.KafkaProducer;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
@Slf4j
public class KafkaProducerImpl implements KafkaProducer {

    private final KafkaTemplate<String, String> kafkaTemplate;

    private final ObjectMapper objectMapper;

    @Override
    public void sendPaymentCreated(BookingInfo bookingInfo) {
        sendMessage(TopicConstant.PAYMENT_CREATED, bookingInfo);
    }

    @Override
    public void sendPaymentCompleted(BookingInfo bookingInfo) {
        sendMessage(TopicConstant.PAYMENT_COMPLETED, bookingInfo);
    }

    @Override
    public void sendPaymentFailed(BookingInfo bookingInfo) {
        sendMessage(TopicConstant.PAYMENT_FAILED, bookingInfo);
    }

    private void sendMessage(String topic, BookingInfo bookingInfo) {
        try {
            String payload = objectMapper.writeValueAsString(bookingInfo);
            kafkaTemplate.send(topic, payload);
            log.debug("Send message to {}, payload: {}", topic, payload);
        } catch (Exception e) {
            log.error("Failed to send payment info to {}, error: {}", topic, e.getMessage());
        }
    }
}
