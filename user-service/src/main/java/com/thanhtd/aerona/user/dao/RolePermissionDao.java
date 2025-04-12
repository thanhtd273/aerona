package com.thanhtd.aerona.user.dao;

import com.thanhtd.aerona.user.model.RolePermission;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface RolePermissionDao extends JpaRepository<RolePermission, String> {
    List<RolePermission> findByRoleId(String roleId);

    List<RolePermission> findByPermissionId(String permissionId);

    RolePermission findByRoleIdAndPermissionId(String roleId, String permissionId);
}
