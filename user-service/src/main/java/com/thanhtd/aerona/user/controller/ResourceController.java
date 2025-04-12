package com.thanhtd.aerona.user.controller;

import com.thanhtd.aerona.user.dto.ResourceInfo;
import com.thanhtd.aerona.user.model.Resource;
import com.thanhtd.aerona.user.service.ResourceService;
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
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RequestBody;

import java.util.List;

@RestController
@RequestMapping("/api/v1/user-service/resources")
@AllArgsConstructor
public class ResourceController {
    private static final Logger logger = LoggerFactory.getLogger(ResourceController.class);

    private final ResourceService resourceService;

    @GetMapping(value = "")
    public APIResponse findAllResources(HttpServletResponse response) {
        long start = System.currentTimeMillis();
        try {
            List<Resource> resources = resourceService.findAll();
            logger.debug("Call API GET /api/v1/user-service/resources success, count = {}", resources.size());
            return new APIResponse(ErrorCode.SUCCESS, "Find all resources success", System.currentTimeMillis() - start, resources);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/resources failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping(value = "/{resourceId}")
    public APIResponse findResourceById(HttpServletResponse response, @PathVariable("resourceId") String resourceId) {
        long start = System.currentTimeMillis();
        try {
            Resource resource = resourceService.findByResourceId(resourceId);
            logger.debug("Call API GET /api/v1/user-service/resources/{} success", resourceId);
            return new APIResponse(ErrorCode.SUCCESS, "Find resource success", System.currentTimeMillis() - start, resource);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/resources/{} failed, error = {}", resourceId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @GetMapping(value = "/")
    public APIResponse findResourceByPath(HttpServletResponse response, @RequestParam(value = "path") String path) {
        long start = System.currentTimeMillis();
        try {
            Resource resource = resourceService.findByPath(path);
            logger.debug("Call API GET /api/v1/user-service/resources?path={} success", path);
            return new APIResponse(ErrorCode.SUCCESS, "Find resource by path successfully", System.currentTimeMillis() - start, resource);
        } catch (Exception e) {
            logger.error("Call API GET /api/v1/user-service/resources?path={} failed, error = {}", path, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PostMapping(value = "")
    public APIResponse createResource(HttpServletResponse response, @RequestBody ResourceInfo resourceInfo) {
        long start = System.currentTimeMillis();
        try {
            Resource resource = resourceService.createResource(resourceInfo);
            logger.debug("Call API POST /api/v1/user-service/resources success");
            return new APIResponse(ErrorCode.SUCCESS, "Create resource successfully", System.currentTimeMillis() - start, resource);
        } catch (Exception e) {
            logger.error("Call API POST /api/v1/user-service/resources failed, error = {}", e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @PutMapping(value = "/{resourceId}")
    public APIResponse updateResource(HttpServletResponse response, @PathVariable("resourceId") String resourceId, @RequestBody ResourceInfo resourceInfo) {
        long start = System.currentTimeMillis();
        try {
            Resource resource = resourceService.updateResource(resourceId, resourceInfo);
            logger.debug("Call API PUT /api/v1/user-service/resources/{} success", resourceId);
            return new APIResponse(ErrorCode.SUCCESS, "Update resource successfully", System.currentTimeMillis() - start, resource);
        } catch (Exception e) {
            logger.error("Call API PUT /api/v1/user-service/resources/{} failed, error = {}", resourceId, e.getMessage());
            return ExceptionHandler.handleException(response, e, start);
        }
    }

    @DeleteMapping(value = "/{resourceId}")
    public APIResponse deleteResource(HttpServletResponse response, @PathVariable("resourceId") String resourceId) {
        long start = System.currentTimeMillis();

        ErrorCode errorCode = resourceService.deleteResource(resourceId);
        logger.debug("Call API DELETE /api/v1/user-service/resources/{} success", resourceId);
        response.setStatus(errorCode.getValue());
        return new APIResponse(errorCode, "Deactivate resource", System.currentTimeMillis() - start, errorCode.getMessage());

    }
}
