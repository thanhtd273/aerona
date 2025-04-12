package com.thanhtd.aerona.user.model;


import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

@Table(name = "permission_resource_action")
@Entity
@Data
public class PermissionResourceAction {
    @Id
    @Column(name = "id")
    private String id;

    @Column(name = "permission_id")
    private String permissionId;

    @Column(name = "resource_id")
    private String resourceId;

    @Column(name = "action_id")
    private String actionId;

    @Column(name = "deleted")
    private Boolean deleted;
}
