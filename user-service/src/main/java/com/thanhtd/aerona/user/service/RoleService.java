package com.thanhtd.aerona.user.service;

import com.thanhtd.aerona.user.dto.RoleInfo;
import com.thanhtd.aerona.user.model.Role;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;

import java.util.List;

public interface RoleService {
    List<Role> findAll();

    Role findByRoleId(String roleId) throws LogicException;

    Role findByName(String name) throws LogicException;

    List<Role> findByUserId(String userId) throws LogicException;

    Role createRole(RoleInfo roleInfo) throws LogicException;

    Role updateRole(String roleId, RoleInfo roleInfo) throws LogicException;

    ErrorCode deleteRole(String roleId);

    ErrorCode addPermission(String roleId, String permissionId) throws LogicException;

    ErrorCode addPermissions(String roleId, List<String> permissionIds) throws LogicException;

    ErrorCode addPermission(RoleInfo roleInfo) throws LogicException;

    ErrorCode removePermission(String roleId, String permissionId) throws LogicException;

    ErrorCode removePermissions(String roleId, List<String> permissionIds) throws LogicException;

    ErrorCode removePermissions(RoleInfo roleInfo) throws LogicException;
}
