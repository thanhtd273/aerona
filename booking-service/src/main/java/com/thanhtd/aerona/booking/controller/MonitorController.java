package com.thanhtd.aerona.booking.controller;

import org.springframework.web.bind.annotation.RestController;

import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;

import org.springframework.web.bind.annotation.GetMapping;


@RestController
public class MonitorController {
    @GetMapping("/health")
    public APIResponse healthCheck() {
        return new APIResponse(ErrorCode.SUCCESS, "Service is healthy", 0L, "Health");
    }
    
}
