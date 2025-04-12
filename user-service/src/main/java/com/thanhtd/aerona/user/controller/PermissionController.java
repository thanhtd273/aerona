package com.thanhtd.aerona.user.controller;

import com.thanhtd.aerona.user.dto.PermissionInfo;
import com.thanhtd.aerona.user.model.Permission;
import com.thanhtd.aerona.user.service.PermissionService;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.core.APIResponse;
import com.thanhtd.aerona.base.exception.ExceptionHandler;
import jakarta.servlet.http.HttpServletResponse;
import lombok.AllArgsConstructor;
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

import java.util.List;

@RestController
@RequestMapping("/api/v1/user-service/permissions")
@AllArgsConstructor
public class PermissionController {
    private static final Logger logger = LoggerFactory.getLogger(PermissionController.class);

    private final PermissionService permissionService;

    @GetMapping(value = "")
    public APIResponse findAllPermissions(HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            List<Permission> permissions = permissionService.findAll();
            logger.debug("Call API GET /api/v1/user-service/permissions success, count = {}", permissions.size());
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, permissions);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/permissions failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping(value = "/{permissionId}")
    public APIResponse findPermissionById(HttpServletResponse response, @PathVariable("permissionId") String permissionId) {
        long start = System.currentTimeMillis();
        try {
            Permission permission = permissionService.findByPermissionId(permissionId);
            logger.debug("Call API GET /api/v1/user-service/permissions/{} success", permissionId);
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, permission);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/permissions/{} failed, error = {}", permissionId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping(value = "")
    public APIResponse createPermission(HttpServletResponse response, @RequestBody PermissionInfo permissionInfo) {
        long start = System.currentTimeMillis();
        try {
            Permission permission = permissionService.createPermission(permissionInfo);
            logger.debug("Call API POST /api/v1/user-service/permissions success");
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, permission);
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/permissions failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PutMapping(value = "/{permissionId}")
    public APIResponse updatePermission(HttpServletResponse response, @PathVariable("permissionId") String permissionId, @RequestBody PermissionInfo permissionInfo) {
        long start = System.currentTimeMillis();
        try {
            Permission permission = permissionService.updatePermission(permissionId, permissionInfo);
            logger.debug("Call API PUT /api/v1/user-service/permissions/{} success", permissionId);
            return new APIResponse(ErrorCode.SUCCESS, "success", System.currentTimeMillis() - start, permission);
        } catch (Exception e) {
            logger.error("Call API PUT /api/v1/user-service/permissions/{} failed, error = {}", permissionId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @DeleteMapping(value = "/{permissionId}")
    public APIResponse deletePermission(HttpServletResponse response, @PathVariable("permissionId") String permissionId) {
        long start = System.currentTimeMillis();

        ErrorCode errorCode = permissionService.deletePermission(permissionId);
        logger.debug("Call API DELETE /api/v1/user-service/permissions/{} success", permissionId);
        return new APIResponse(errorCode, errorCode.getMessage(), System.currentTimeMillis() - start, errorCode.getMessage());

    }
}
