package com.thanhtd.aerona.user.dao;

import com.thanhtd.aerona.user.model.PermissionResourceAction;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface PermissionResourceActionDao extends JpaRepository<PermissionResourceAction, String> {
    List<PermissionResourceAction> findByPermissionId(String permissionId);


    List<PermissionResourceAction> findByResourceId(String resourceId);

}
