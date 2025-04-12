package com.thanhtd.aerona.user.service;

import com.thanhtd.aerona.user.dto.PermissionInfo;
import com.thanhtd.aerona.user.model.Permission;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;

import java.util.List;

public interface PermissionService {
    List<Permission> findAll();

    Permission findByPermissionId(String permissionId) throws LogicException;

    List<Permission> findByRoleId(String roleId) throws LogicException;

    Permission createPermission(PermissionInfo permissionInfo) throws LogicException;

    Permission updatePermission(String permissionId, PermissionInfo permissionInfo) throws LogicException;

    ErrorCode deletePermission(String permissionId);
}
