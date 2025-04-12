package com.thanhtd.aerona.user.service.impl;

import com.thanhtd.aerona.user.dao.PermissionDao;
import com.thanhtd.aerona.user.dao.RolePermissionDao;
import com.thanhtd.aerona.user.dto.PermissionInfo;
import com.thanhtd.aerona.user.model.Permission;
import com.thanhtd.aerona.user.model.RolePermission;
import com.thanhtd.aerona.user.service.PermissionService;
import com.thanhtd.aerona.base.constant.DataStatus;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.util.ObjectUtils;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

@Service
@RequiredArgsConstructor
public class PermissionServiceImpl implements PermissionService {

    private final PermissionDao permissionDao;

    private final RolePermissionDao rolePermissionDao;

    @Override
    public List<Permission> findAll() {
        return permissionDao.findAll();
    }

    @Override
    public Permission findByPermissionId(String permissionId) throws LogicException {
        if (ObjectUtils.isEmpty(permissionId))
            throw new LogicException(ErrorCode.ID_NULL);

        return permissionDao.findByPermissionId(permissionId);
    }

    @Override
    public List<Permission> findByRoleId(String roleId) throws LogicException {
        List<Permission> permissions = new ArrayList<>();
        List<RolePermission> rolePermissions = rolePermissionDao.findByRoleId(roleId);
        for (RolePermission rolePermission : rolePermissions) {
            Permission permission = permissionDao.findByPermissionId(rolePermission.getPermissionId());
            permissions.add(permission);
        }
        return permissions;
    }

    @Override
    public Permission createPermission(PermissionInfo permissionInfo) throws LogicException {
        if (ObjectUtils.isEmpty(permissionInfo))
            throw new LogicException(ErrorCode.DATA_NULL);
        if (ObjectUtils.isEmpty(permissionInfo.getName()))
            throw new LogicException(ErrorCode.BLANK_FIELD);

        Permission permission = new Permission();
        permission.setName(permissionInfo.getName());
        if (!ObjectUtils.isEmpty(permissionInfo.getDescription()))
            permission.setDescription(permissionInfo.getDescription());
        if (!ObjectUtils.isEmpty(permissionInfo.getStatus()))
            permission.setStatus(permissionInfo.getStatus());

        permission.setCreateDate(new Date());
        permission.setStatus(DataStatus.ACTIVE);
        permission = permissionDao.save(permission);

        return permission;
    }

    @Override
    public Permission updatePermission(String permissionId, PermissionInfo permissionInfo) throws LogicException {
        if (ObjectUtils.isEmpty(permissionInfo) || ObjectUtils.isEmpty(permissionId))
            throw new LogicException(ErrorCode.DATA_NULL);

        Permission permission = permissionDao.findByPermissionId(permissionId);
        if (ObjectUtils.isEmpty(permission))
            throw new LogicException(ErrorCode.DATA_NULL);

        if (!ObjectUtils.isEmpty(permissionInfo.getName()))
            permission.setName(permissionInfo.getName());
        if (!ObjectUtils.isEmpty(permissionInfo.getDescription())) {
            permission.setDescription(permissionInfo.getDescription());
        }
        if (!ObjectUtils.isEmpty(permissionInfo.getStatus())) {
            permission.setStatus(permissionInfo.getStatus());
        }
        permission.setModifiedDate(new Date());
        permission = permissionDao.save(permission);

        return permission;
    }

    @Override
    public ErrorCode deletePermission(String permissionId) {
        if (ObjectUtils.isEmpty(permissionId))
            return ErrorCode.ID_NULL;
        Permission permission = permissionDao.findByPermissionId(permissionId);
        if (ObjectUtils.isEmpty(permission))
            return ErrorCode.DATA_NULL;
        permission.setModifiedDate(new Date());
        permission.setStatus(DataStatus.DELETED);
        permissionDao.save(permission);
        return ErrorCode.SUCCESS;
    }
}
