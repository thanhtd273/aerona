//package com.thanhtd.aerona.booking.controller;
//
//import com.thanhtd.aerona.base.constant.ErrorCode;
//import com.thanhtd.aerona.base.core.APIResponse;
//import com.thanhtd.aerona.base.exception.ExceptionHandler;
//import com.thanhtd.aerona.booking.dto.BookingInfo;
//import com.thanhtd.aerona.booking.service.KafkaProducer;
//import jakarta.servlet.http.HttpServletResponse;
//import lombok.RequiredArgsConstructor;
//import lombok.extern.slf4j.Slf4j;
//
//import java.util.UUID;
//
//import org.springframework.web.bind.annotation.GetMapping;
//import org.springframework.web.bind.annotation.RequestMapping;
//import org.springframework.web.bind.annotation.RestController;
//
//@RestController
//@RequestMapping("/test")
//@RequiredArgsConstructor
//@Slf4j
//public class TestController {
//
//    private final KafkaProducer kafkaProducer;
//
//    @GetMapping()
//    public APIResponse test1(HttpServletResponse response) {
//        long start = System.currentTimeMillis();
//        try {
//            BookingInfo bookingInfo = new BookingInfo();
//           bookingInfo.setBookingId(UUID.randomUUID().toString());
//           bookingInfo.setStatus("Test");
//            kafkaProducer.sendBookingCreated(bookingInfo);
//            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, bookingInfo);
//        } catch (Exception e) {
//            log.error("Test error: {}", e.getMessage());
//            return ExceptionHandler.handleException(response, e, start);
//        }
//    }
//}
