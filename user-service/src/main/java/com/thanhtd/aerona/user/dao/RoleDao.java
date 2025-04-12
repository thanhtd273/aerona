package com.thanhtd.aerona.user.dao;

import com.thanhtd.aerona.user.model.Role;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface RoleDao extends JpaRepository<Role, String> {

    @Query("SELECT u FROM Role u WHERE u.status <> 0")
    List<Role> findAll();

    @Query("SELECT u FROM Role u WHERE u.status <> 0 AND u.roleId = :roleId")
    Role findByRoleId(@Param(value = "roleId") String roleId);

    @Query("SELECT u FROM Role u WHERE u.status <> 0 AND u.name = :name")
    Role findByName(@Param(value = "name") String name);
}
