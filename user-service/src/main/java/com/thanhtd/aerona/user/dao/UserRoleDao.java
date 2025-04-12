package com.thanhtd.aerona.user.dao;

import com.thanhtd.aerona.user.model.UserRole;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface UserRoleDao extends JpaRepository<UserRole, String> {
    List<UserRole> findByUserId(String userId);

    List<UserRole> findByRoleId(String roleId);

    UserRole findByUserIdAndRoleId(String userId, String roleId);
}
