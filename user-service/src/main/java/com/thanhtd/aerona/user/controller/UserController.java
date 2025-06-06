package com.thanhtd.aerona.user.controller;

import com.thanhtd.aerona.user.dto.OTPInfo;
import com.thanhtd.aerona.user.dto.UpdatePasswordInfo;
import com.thanhtd.aerona.user.dto.UserInfo;
import com.thanhtd.aerona.user.model.User;
import com.thanhtd.aerona.user.service.UserService;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;
import com.thanhtd.aerona.base.exception.ExceptionHandler;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;


@RestController
@RequestMapping("/api/v1/user-service")
@RequiredArgsConstructor
public class UserController {
    private static final Logger logger = LoggerFactory.getLogger(UserController.class);

    private final UserService userService;

    @GetMapping("/users/{userId}")
    public APIResponse findByUserId(HttpServletResponse response, @PathVariable("userId") String userId) {
        long start = System.currentTimeMillis();
        try {
            User user = userService.findByUserId(userId);
            logger.debug("Call API GET /api/v1/user-service/users/{} success", userId);
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, user);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/users/{} failed, error: {}", userId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping("/users/me")
    public APIResponse getCurrentUser(HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            UserInfo userInfo = userService.getCurrentUser();
            logger.debug("Call API /api/v1/user-service/users/me success, took: {}", System.currentTimeMillis() - start);
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, userInfo);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/users/me failed, error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PutMapping("/users/{userId}")
    public APIResponse updateUser(HttpServletResponse response, @PathVariable("userId") String userId, @RequestBody UserInfo userInfo) {
        long start = System.currentTimeMillis();
        try {
            User user = userService.updateUser(userId, userInfo);
            logger.debug("Call API PUT /api/v1/user-service/users/{} success", userId);
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, user);
        } catch (Exception e) {
            logger.error("Call API PUT /api/v1/user-service/users/{} failed, request body: {}, error: {}", userId, userInfo, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @DeleteMapping("/users/{userId}")
    public APIResponse deactivateUser(HttpServletResponse response, @PathVariable("userId") String userId) {
        long start = System.currentTimeMillis();
        try {
            ErrorCode errorCode = userService.deactivateUser(userId);
            logger.debug("Call API DELETE /api/v1/user-service/users/{} success", userId);
            response.setStatus(errorCode.getValue());
            return new APIResponse(errorCode, "", System.currentTimeMillis() - start, errorCode.getMessage());
        } catch (Exception e) {
            logger.error("Call API DELETE /api/v1/user-service/users/{} failed, error: {}", userId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping("/users/active")
    public APIResponse activateUser(HttpServletResponse response, @RequestBody OTPInfo activationInfo) {
        long start = System.currentTimeMillis();
        try {
            ErrorCode errorCode = userService.activateUser(activationInfo);
            if (errorCode == ErrorCode.SUCCESS) {
                logger.info("Call API POST /api/v1/user-service/users/activate success, email: {}", activationInfo.getEmail());
            } else {
                logger.error("Call API POST /api/v1/user-service/users/activate failed, email: {}, error: {}",
                        activationInfo.getEmail(), errorCode.getMessage());
            }
            response.setStatus(errorCode.getValue());
            return new APIResponse(errorCode, "", System.currentTimeMillis() - start, errorCode.getMessage());
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/users/active failed, user email: {}, error: {}",
                    activationInfo.getEmail(), e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping("/users/change-password")
    public APIResponse changePassword(HttpServletResponse response, @RequestBody UpdatePasswordInfo updatePasswordInfo) {
        long start = System.currentTimeMillis();
        try {
            ErrorCode errorCode = userService.changePassword(updatePasswordInfo);
            if (errorCode == ErrorCode.SUCCESS) {
                logger.info("Call API POST /api/v1/user-service/users/change-password success, user email: {}", updatePasswordInfo.getEmail());
            } else {
                logger.error("Call API POST /api/v1/user-service/users/change-password failed, user email: {}, error: {}",
                        updatePasswordInfo.getEmail(), errorCode.getMessage());
            }

            response.setStatus(errorCode.getValue());
            return new APIResponse(errorCode, "", System.currentTimeMillis() - start, errorCode.getMessage());
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/users/change-password failed, error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping("/users/forgot-password")
    public APIResponse forgotPassword(HttpServletResponse response, @RequestBody UserInfo userInfo) {
        long start = System.currentTimeMillis();
        try {
            ErrorCode errorCode = userService.forgotPassword(userInfo);
            if (errorCode == ErrorCode.SUCCESS) {
                logger.info("Call API POST /api/v1/user-service/users/forgot-password success, user email: {}", userInfo.getEmail());
            } else {
                logger.error("Call API POST /api/v1/user-service/users/forgot-password failed, user email: {}, error: ", userInfo.getEmail());
            }

            response.setStatus(errorCode.getValue());
            return new APIResponse(errorCode, "", System.currentTimeMillis() - start, errorCode.getMessage());
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/users/forgot-password failed, user email: {}, error: {}",
                    userInfo.getEmail(), e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping("/users/verify-password-reset-code")
    public APIResponse verifyPasswordResetCode(HttpServletResponse response, @RequestBody OTPInfo passwordResetInfo) {
        long start = System.currentTimeMillis();
        try {
            ErrorCode errorCode = userService.verifyPasswordResetCode(passwordResetInfo);
            if (errorCode == ErrorCode.SUCCESS) {
                logger.info("Call API POST /api/v1/user-service/users/verify-password-reset-code success, user email: {}", passwordResetInfo.getEmail());
            } else {
                logger.error("Call API POST /api/v1/user-service/users/verify-password-reset-code failed, user email: {}, error: {}",
                        passwordResetInfo.getEmail(), errorCode.getMessage());
            }

            response.setStatus(errorCode.getValue());
            return new APIResponse(errorCode, "", System.currentTimeMillis() - start, errorCode.getMessage());
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/users/verify-password-reset-code failed, error: {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }


}
