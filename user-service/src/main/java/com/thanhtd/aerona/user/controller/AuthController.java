package com.thanhtd.aerona.user.controller;

import com.thanhtd.aerona.user.dto.Credential;
import com.thanhtd.aerona.user.dto.LoginSessionInfo;
import com.thanhtd.aerona.user.dto.UserInfo;
import com.thanhtd.aerona.user.model.User;
import com.thanhtd.aerona.user.service.AuthService;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;
import com.thanhtd.aerona.base.exception.ExceptionHandler;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpHeaders;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping(value = "/api/v1/user-service")
@RequiredArgsConstructor
public class AuthController {
    private static final Logger logger = LoggerFactory.getLogger(AuthController.class);

    private final AuthService authService;

    @GetMapping(value = "/verify-token")
    public APIResponse verifyToken(@RequestHeader(value = HttpHeaders.AUTHORIZATION) String token, HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {

            logger.info("Call API GET /api/v1/user-service/verify-token success");
            UserInfo userInfo = authService.verifyToken(token);
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, userInfo);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/verify-token fail, error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping(value = "/login")
    public APIResponse login(HttpServletRequest request, HttpServletResponse response, @RequestBody Credential credential) {
        long start = System.currentTimeMillis();
        try {
            LoginSessionInfo loginSessionInfo = authService.login(request, credential);
            logger.info("Call API /api/v1/user-service/login success");
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, loginSessionInfo);
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/login fail, error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping("/sign-up")
    public APIResponse createUser(HttpServletResponse response, @RequestBody UserInfo userInfo) {
        long start = System.currentTimeMillis();
        try {
            User user = authService.signUp(userInfo);
            logger.debug("Call API POST /api/v1/user-service/sign-up success, userId = {}", user.getUserId());
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, user);
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/sign-up failed, request body: {}, error: {}", userInfo, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }
}
