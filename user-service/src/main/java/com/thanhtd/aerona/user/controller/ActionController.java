package com.thanhtd.aerona.user.controller;

import com.thanhtd.aerona.user.dto.ActionInfo;
import com.thanhtd.aerona.user.model.Action;
import com.thanhtd.aerona.user.service.ActionService;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;
import com.thanhtd.aerona.base.exception.ExceptionHandler;
import jakarta.servlet.http.HttpServletResponse;
import lombok.AllArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/user-service/actions")
@AllArgsConstructor
public class ActionController {
    private static final Logger logger = LoggerFactory.getLogger(ActionController.class);

    private final ActionService actionService;

    @GetMapping(value = "")
    public APIResponse findAllActions(HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            List<Action> actions = actionService.findAll();
            logger.debug("Call API GET /api/v1/user-service/actions success, count = {}", actions.size());
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, actions);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/actions failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping(value = "/{actionId}")
    public APIResponse findActionById(HttpServletResponse response, @PathVariable("actionId") String actionId) {
        long start = System.currentTimeMillis();
        try {
            Action action = actionService.findByActionId(actionId);
            logger.debug("Call API GET /api/v1/user-service/actions/{} success", actionId);
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, action);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/actions/{} failed, error = {}", actionId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping(value = "/")
    public APIResponse findActionByCode(HttpServletResponse response, @RequestParam("code") Integer code) {
        long start = System.currentTimeMillis();
        try {
            Action action = actionService.findByCode(code);
            logger.debug("Call API GET /api/v1/user-service/actions?code={} success", code);
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, action);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/actions?code={} failed, error = {}", code, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping(value = "")
    public APIResponse createAction(HttpServletResponse response, @RequestBody ActionInfo actionInfo) {
        long start = System.currentTimeMillis();
        try {
            Action action = actionService.createAction(actionInfo);
            logger.debug("Call API POST /api/v1/user-service/actions success");
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, action);
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/actions failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PutMapping(value = "/{actionId}")
    public APIResponse updateAction(HttpServletResponse response, @PathVariable("actionId") String actionId, @RequestBody ActionInfo actionInfo) {
        long start = System.currentTimeMillis();
        try {
            Action action = actionService.updateAction(actionId, actionInfo);
            logger.debug("Call API PUT /api/v1/user-service/actions/{} success", actionId);
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, action);
        } catch (Exception e) {
            logger.error("Call API PUT /api/v1/user-service/actions/{} failed, error = {}", actionId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @DeleteMapping(value = "/{actionId}")
    public APIResponse deleteAction(HttpServletResponse response, @PathVariable("actionId") String actionId) {
        long start = System.currentTimeMillis();

        ErrorCode errorCode = actionService.deleteAction(actionId);
        logger.debug("Call API DELETE /api/v1/user-service/actions/{} success", actionId);
        return new APIResponse(errorCode, errorCode.getMessage(), System.currentTimeMillis() - start, errorCode.getMessage());

    }
}
