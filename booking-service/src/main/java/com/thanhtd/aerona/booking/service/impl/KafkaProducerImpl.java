package com.thanhtd.aerona.booking.service.impl;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.thanhtd.aerona.booking.constants.TopicConstant;
import com.thanhtd.aerona.booking.dto.BookingInfo;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import com.thanhtd.aerona.booking.service.KafkaProducer;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
@Slf4j
public class KafkaProducerImpl implements KafkaProducer {

    private final KafkaTemplate<String, String> kafkaTemplate;

    private final ObjectMapper objectMapper;

    @Override
    public void sendPaymentCompleted(BookingInfo bookingInfo) {
        try {
            String message = objectMapper.writeValueAsString(bookingInfo);
            kafkaTemplate.send(TopicConstant.PAYMENT_COMPLETED, message);
        } catch (JsonProcessingException e) {
            log.error("Failed to send booking info to {}, error: {}", TopicConstant.PAYMENT_COMPLETED, e.getMessage());
        }
    }
}
