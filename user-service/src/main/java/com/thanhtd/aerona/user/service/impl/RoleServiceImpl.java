package com.thanhtd.aerona.user.service.impl;

import com.thanhtd.aerona.user.dao.PermissionDao;
import com.thanhtd.aerona.user.dao.RoleDao;
import com.thanhtd.aerona.user.dao.RolePermissionDao;
import com.thanhtd.aerona.user.dao.UserRoleDao;
import com.thanhtd.aerona.user.dto.RoleInfo;
import com.thanhtd.aerona.user.model.Permission;
import com.thanhtd.aerona.user.model.Role;
import com.thanhtd.aerona.user.model.RolePermission;
import com.thanhtd.aerona.user.model.UserRole;
import com.thanhtd.aerona.user.service.RoleService;
import com.thanhtd.aerona.base.constant.DataStatus;
import com.thanhtd.aerona.base.constant.ErrorCode;
import com.thanhtd.aerona.base.exception.LogicException;
import jodd.util.StringPool;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.ObjectUtils;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

@Service
@RequiredArgsConstructor
public class RoleServiceImpl implements RoleService {

    private final RoleDao roleDao;

    private final UserRoleDao userRoleDao;

    private final RolePermissionDao rolePermissionDao;

    private final PermissionDao permissionDao;

    @Override
    public List<Role> findAll() {
        return roleDao.findAll();
    }

    @Override
    public Role findByRoleId(String roleId) throws LogicException {
        if (ObjectUtils.isEmpty(roleId))
            throw new LogicException(ErrorCode.ID_NULL);

        return roleDao.findByRoleId(roleId);
    }

    @Override
    public Role findByName(String name) throws LogicException {
        return roleDao.findByName(name);
    }

    @Override
    public List<Role> findByUserId(String userId) throws LogicException {
        List<Role> roles = new ArrayList<>();
        List<UserRole> userRoles = userRoleDao.findByUserId(userId);
        for (UserRole userRole : userRoles) {
            Role role = findByRoleId(userRole.getRoleId());
            roles.add(role);
        }
        return roles;
    }

    @Override
    public Role createRole(RoleInfo roleInfo) throws LogicException {
        if (ObjectUtils.isEmpty(roleInfo))
            throw new LogicException(ErrorCode.DATA_NULL);
        if (ObjectUtils.isEmpty(roleInfo.getName()))
            throw new LogicException(ErrorCode.BLANK_FIELD);

        Role role = new Role();
        role.setName(roleInfo.getName());
        if (!ObjectUtils.isEmpty(roleInfo.getDescription()))
            role.setDescription(roleInfo.getDescription());
        if (!ObjectUtils.isEmpty(roleInfo.getStatus()))
            role.setStatus(roleInfo.getStatus());

        role.setCreateDate(new Date());
        role.setStatus(DataStatus.ACTIVE);
        role = roleDao.save(role);

        return role;
    }

    @Override
    public Role updateRole(String roleId, RoleInfo roleInfo) throws LogicException {
        if (ObjectUtils.isEmpty(roleInfo) || ObjectUtils.isEmpty(roleId))
            throw new LogicException(ErrorCode.DATA_NULL);

        Role role = roleDao.findByRoleId(roleId);
        if (ObjectUtils.isEmpty(role))
            throw new LogicException(ErrorCode.DATA_NULL);

        if (!ObjectUtils.isEmpty(roleInfo.getName()))
            role.setName(roleInfo.getName());
        if (!ObjectUtils.isEmpty(roleInfo.getDescription())) {
            role.setDescription(roleInfo.getDescription());
        }
        if (!ObjectUtils.isEmpty(roleInfo.getStatus())) {
            role.setStatus(roleInfo.getStatus());
        }
        role.setModifiedDate(new Date());
        role = roleDao.save(role);

        return role;
    }

    @Override
    public ErrorCode deleteRole(String roleId) {
        if (ObjectUtils.isEmpty(roleId))
            return ErrorCode.ID_NULL;
        Role role = roleDao.findByRoleId(roleId);
        if (ObjectUtils.isEmpty(role))
            return ErrorCode.DATA_NULL;
        role.setModifiedDate(new Date());
        role.setStatus(DataStatus.DELETED);
        return ErrorCode.SUCCESS;
    }

    @Override
    public ErrorCode addPermission(String roleId, String permissionId) throws LogicException {
        Role role = findByRoleId(roleId);
        if (ObjectUtils.isEmpty(role))
            return ErrorCode.NOT_FOUND_ROLE;
        Permission permission = permissionDao.findByPermissionId(permissionId);
        if (ObjectUtils.isEmpty(permission))
            return ErrorCode.NOT_FOUND_PERMISSION;
        RolePermission rolePermission = new RolePermission();
        rolePermission.setRoleId(roleId);
        rolePermission.setPermissionId(permissionId);

        rolePermissionDao.save(rolePermission);
        return ErrorCode.SUCCESS;
    }

    @Override
    public ErrorCode addPermissions(String roleId, List<String> permissionIds) throws LogicException {
        for (String permissionId : permissionIds) {
            addPermission(roleId, permissionId);
        }
        return ErrorCode.SUCCESS;
    }

    @Override
    @Transactional
    public ErrorCode addPermission(RoleInfo roleInfo) throws LogicException {
        if (ObjectUtils.isEmpty(roleInfo))
            return ErrorCode.DATA_NULL;
        String roleId = roleInfo.getRoleId();
        String permissionsStr = roleInfo.getPermissions();
        if (ObjectUtils.isEmpty(permissionsStr))
            return ErrorCode.BLANK_FIELD;

        String[] permissionsNameArr = permissionsStr.split(StringPool.DASH);
        for (String permissionName : permissionsNameArr) {
            Permission permission = permissionDao.findByName(permissionName);
            if (ObjectUtils.isEmpty(permission)) continue;
            addPermission(roleId, permission.getPermissionId());
        }

        return ErrorCode.SUCCESS;
    }

    @Override
    public ErrorCode removePermission(String roleId, String permissionId) throws LogicException {
        RolePermission rolePermission = rolePermissionDao.findByRoleIdAndPermissionId(roleId, permissionId);
        if (ObjectUtils.isEmpty(rolePermission))
            return ErrorCode.DATA_NULL;
        rolePermissionDao.delete(rolePermission);
        return ErrorCode.SUCCESS;
    }

    @Override
    public ErrorCode removePermissions(String roleId, List<String> permissionIds) throws LogicException {
        for (String permissionId : permissionIds) {
            removePermission(roleId, permissionId);
        }
        return ErrorCode.SUCCESS;
    }

    @Override
    public ErrorCode removePermissions(RoleInfo roleInfo) throws LogicException {
        if (ObjectUtils.isEmpty(roleInfo))
            return ErrorCode.DATA_NULL;
        String roleId = roleInfo.getRoleId();
        String permissionsStr = roleInfo.getPermissions();
        if (ObjectUtils.isEmpty(permissionsStr))
            return ErrorCode.BLANK_FIELD;

        String[] permissionsNameArr = permissionsStr.split(StringPool.DASH);
        for (String permissionName : permissionsNameArr) {
            Permission permission = permissionDao.findByName(permissionName);
            if (ObjectUtils.isEmpty(permission)) continue;
            removePermission(roleId, permission.getPermissionId());
        }

        return ErrorCode.SUCCESS;
    }
}
