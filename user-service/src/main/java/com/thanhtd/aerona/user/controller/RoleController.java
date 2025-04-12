package com.thanhtd.aerona.user.controller;

import com.thanhtd.aerona.user.dto.RoleInfo;
import com.thanhtd.aerona.user.model.Role;
import com.thanhtd.aerona.user.service.RoleService;
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
@RequestMapping("/api/v1/user-service/roles")
@AllArgsConstructor
public class RoleController {
    private static final Logger logger = LoggerFactory.getLogger(RoleController.class);

    private final RoleService roleService;

    @GetMapping(value = "")
    public APIResponse findAllRoles(HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            List<Role> roles = roleService.findAll();
            logger.debug("Call API GET /api/v1/user-service/roles success, count = {}", roles.size());
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, roles);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/roles failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping(value = "/{roleId}")
    public APIResponse findRoleById(HttpServletResponse response, @PathVariable("roleId") String roleId) {
        long start = System.currentTimeMillis();
        try {
            Role role = roleService.findByRoleId(roleId);
            logger.debug("Call API GET /api/v1/user-service/roles/{} success", roleId);
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, role);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/roles/{} failed, error = {}", roleId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping(value = "")
    public APIResponse createRole(HttpServletResponse response, @RequestBody RoleInfo roleInfo) {
        long start = System.currentTimeMillis();
        try {
            Role role = roleService.createRole(roleInfo);
            logger.debug("Call API POST /api/v1/user-service/roles success");
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, role);
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/roles failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PutMapping(value = "/{roleId}")
    public APIResponse updateRole(HttpServletResponse response, @PathVariable("roleId") String roleId, @RequestBody RoleInfo roleInfo) {
        long start = System.currentTimeMillis();
        try {
            Role role = roleService.updateRole(roleId, roleInfo);
            logger.debug("Call API PUT /api/v1/user-service/roles/{} success", roleId);
            return new APIResponse(ErrorCode.SUCCESS, "", System.currentTimeMillis() - start, role);
        } catch (Exception e) {
            logger.error("Call API PUT /api/v1/user-service/roles/{} failed, error = {}", roleId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @DeleteMapping(value = "/{roleId}")
    public APIResponse deleteRole(HttpServletResponse response, @PathVariable("roleId") String roleId) {
        long start = System.currentTimeMillis();

        ErrorCode errorCode = roleService.deleteRole(roleId);
        logger.debug("Call API DELETE /api/v1/user-service/roles/{} success", roleId);
        response.setStatus(errorCode.getValue());
        return new APIResponse(errorCode, "Deactivate role success", System.currentTimeMillis() - start, errorCode.getMessage());
    }

    @PostMapping("/add-permission")
    public APIResponse addPermission(HttpServletResponse response, @RequestBody RoleInfo roleInfo ) {
        long start = System.currentTimeMillis();
        try {
            ErrorCode errorCode = roleService.addPermission(roleInfo);
            if (errorCode.equals(ErrorCode.SUCCESS)) {
                logger.debug("Call API POST /api/v1/user-service/roles/add-permission success");
            } else {
                logger.error("Call API POST /api/v1/user-service/roles/add-permission failed, error = {}", errorCode.getMessage());
            }
            response.setStatus(errorCode.getValue());
            return new APIResponse(errorCode, "Add permissions successfully", System.currentTimeMillis() - start, errorCode.getMessage());
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/roles/add-permission failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping("/remove-permission")
    public APIResponse removePermissions(HttpServletResponse response, @RequestBody RoleInfo roleInfo ) {
        long start = System.currentTimeMillis();
        try {
            ErrorCode errorCode = roleService.removePermissions(roleInfo);
            if (errorCode.equals(ErrorCode.SUCCESS)) {
                logger.debug("Call API POST /api/v1/user-service/roles/remove-permission success");
            } else {
                logger.error("Call API POST /api/v1/user-service/roles/remove-permission failed, error = {}", errorCode.getMessage());
            }
            response.setStatus(errorCode.getValue());
            return new APIResponse(errorCode, "Remove permissions success", System.currentTimeMillis() - start, errorCode.getMessage());
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/roles/remove-permission failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }
}
