package com.thanhtd.aerona.user.service.impl;

import com.thanhtd.aerona.user.dao.ResourceDao;
import com.thanhtd.aerona.user.dto.ResourceInfo;
import com.thanhtd.aerona.user.model.Resource;
import com.thanhtd.aerona.user.service.ResourceService;
import com.thanhtd.aerona.base.constant.DataStatus;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.util.ObjectUtils;

import java.util.Date;
import java.util.List;

@Service
@RequiredArgsConstructor
public class ResourceServiceImpl implements ResourceService {
    private final ResourceDao resourceDao;

    @Override
    public List<Resource> findAll() {
        return resourceDao.findAll();
    }

    @Override
    public Resource findByResourceId(String resourceId) throws LogicException {
        if (ObjectUtils.isEmpty(resourceId))
            throw new LogicException(ErrorCode.ID_NULL);

        return resourceDao.findByResourceId(resourceId);
    }

    @Override
    public Resource findByPath(String path) throws LogicException {
        if (ObjectUtils.isEmpty(path))
            throw new LogicException(ErrorCode.NULL_VALUE);

        return resourceDao.findByPath(path);
    }

    @Override
    public Resource createResource(ResourceInfo resourceInfo) throws LogicException {
        if (ObjectUtils.isEmpty(resourceInfo))
            throw new LogicException(ErrorCode.DATA_NULL);
        if (ObjectUtils.isEmpty(resourceInfo.getName()))
            throw new LogicException(ErrorCode.BLANK_FIELD);

        Resource resource = new Resource();
        resource.setName(resourceInfo.getName());
        if (!ObjectUtils.isEmpty(resourceInfo.getDescription()))
            resource.setDescription(resourceInfo.getDescription());
        if (!ObjectUtils.isEmpty(resourceInfo.getStatus()))
            resource.setStatus(resourceInfo.getStatus());

        resource.setCreateDate(new Date());
        resource.setStatus(DataStatus.ACTIVE);
        resource = resourceDao.save(resource);

        return resource;
    }

    @Override
    public Resource updateResource(String resourceId, ResourceInfo resourceInfo) throws LogicException {
        if (ObjectUtils.isEmpty(resourceInfo) || ObjectUtils.isEmpty(resourceId))
            throw new LogicException(ErrorCode.DATA_NULL);

        Resource resource = resourceDao.findByResourceId(resourceId);
        if (ObjectUtils.isEmpty(resource))
            throw new LogicException(ErrorCode.DATA_NULL);

        if (!ObjectUtils.isEmpty(resourceInfo.getName()))
            resource.setName(resourceInfo.getName());
        if (!ObjectUtils.isEmpty(resourceInfo.getDescription())) {
            resource.setDescription(resourceInfo.getDescription());
        }
        if (!ObjectUtils.isEmpty(resourceInfo.getStatus())) {
            resource.setStatus(resourceInfo.getStatus());
        }
        resource.setModifiedDate(new Date());
        resource = resourceDao.save(resource);

        return resource;
    }

    @Override
    public ErrorCode deleteResource(String resourceId) {
        if (ObjectUtils.isEmpty(resourceId))
            return ErrorCode.ID_NULL;
        Resource resource = resourceDao.findByResourceId(resourceId);
        if (ObjectUtils.isEmpty(resource))
            return ErrorCode.DATA_NULL;
        resource.setModifiedDate(new Date());
        resource.setStatus(DataStatus.DELETED);
        return ErrorCode.SUCCESS;
    }
}
