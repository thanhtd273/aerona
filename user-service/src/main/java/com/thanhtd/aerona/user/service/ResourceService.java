package com.thanhtd.aerona.user.service;

import com.thanhtd.aerona.user.dto.ResourceInfo;
import com.thanhtd.aerona.user.model.Resource;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;

import java.util.List;

public interface ResourceService {
    List<Resource> findAll();

    Resource findByResourceId(String resourceId) throws LogicException;

    Resource findByPath(String path) throws LogicException;

    Resource createResource(ResourceInfo resourceInfo) throws LogicException;

    Resource updateResource(String resourceId, ResourceInfo resourceInfo) throws LogicException;

    ErrorCode deleteResource(String resourceId);
}
