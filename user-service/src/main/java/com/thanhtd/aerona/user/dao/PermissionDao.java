package com.thanhtd.aerona.user.dao;

import com.thanhtd.aerona.user.model.Permission;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface PermissionDao extends JpaRepository<Permission, String> {

    @Query("SELECT u FROM Permission u WHERE u.status <> 0")
    List<Permission> findAll();

    @Query("SELECT u FROM Permission u WHERE u.status <> 0 AND u.permissionId = :permissionId")
    Permission findByPermissionId(@Param("permissionId") String permissionId);

    @Query("SELECT u FROM Permission u WHERE u.status <> 0 AND u.name = :name")
    Permission findByName(@Param("name") String name);
}
