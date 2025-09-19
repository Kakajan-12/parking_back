package repository_test

import (
	"testing"

	"backend/contrib/models"
	"backend/contrib/repository"
	"backend/test"
)

func TestCameraRepository_CRUD(t *testing.T) {
	db := test.TestDB
	cameraRepo := repository.NewCameraRepository(db)

	// Create
	camera, err := cameraRepo.CameraCreate("Test Camera", models.CameraType("IP"))
	if err != nil {
		t.Fatalf("CameraCreate failed: %v", err)
	}
	if camera.ID == 0 {
		t.Errorf("Expected camera ID to be set")
	}

	// Get by ID
	gotCamera, err := cameraRepo.CameraGetByID(camera.ID, false)
	if err != nil {
		t.Fatalf("CameraGetByID failed: %v", err)
	}
	if gotCamera.Name != "Test Camera" {
		t.Errorf("Expected Name 'Test Camera', got '%s'", gotCamera.Name)
	}

	// List
	cameras, err := cameraRepo.CameraList(nil)
	if err != nil {
		t.Fatalf("CameraList failed: %v", err)
	}
	if len(cameras) != 1 {
		t.Errorf("Expected 1 camera, got %d", len(cameras))
	}

	// Count
	count, err := cameraRepo.CameraCount(nil)
	if err != nil {
		t.Fatalf("CameraCount failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Update
	gotCamera.Name = "Updated Camera"
	err = cameraRepo.CameraUpdate(gotCamera)
	if err != nil {
		t.Fatalf("CameraUpdate failed: %v", err)
	}
	updatedCamera, _ := cameraRepo.CameraGetByID(camera.ID, false)
	if updatedCamera.Name != "Updated Camera" {
		t.Errorf("Expected updated name, got '%s'", updatedCamera.Name)
	}

	// Delete (soft delete)
	err = cameraRepo.CameraDelete(camera.ID)
	if err != nil {
		t.Fatalf("CameraDelete failed: %v", err)
	}

	deletedCamera, _ := cameraRepo.CameraGetByID(camera.ID, false)
	if deletedCamera != nil {
		t.Errorf("Expected camera to be soft-deleted")
	}

	// Include deleted
	includedDeletedCamera, _ := cameraRepo.CameraGetByID(camera.ID, true)
	if includedDeletedCamera == nil {
		t.Errorf("Expected to retrieve soft-deleted camera when includeDeleted=true")
	}
	if includedDeletedCamera.DeletedAt == nil {
		t.Errorf("Expected DeletedAt to be set")
	}
	if includedDeletedCamera.Name != "Updated Camera" {
		t.Errorf("Expected Name 'Updated Camera', got '%s'", includedDeletedCamera.Name)
	}
}

func TestCameraRepository_SearchAndFilter(t *testing.T) {
	db := test.TestDB
	cameraRepo := repository.NewCameraRepository(db)

	// Seed cameras
	cameraRepo.CameraCreate("Front Door", models.CameraType("IP"))
	cameraRepo.CameraCreate("Back Door", models.CameraType("Analog"))
	cameraRepo.CameraCreate("Garage", models.CameraType("IP"))

	// Search
	searchTerm := "Door"
	cameras, err := cameraRepo.CameraList(&repository.CameraRepoFilter{
		Search: &searchTerm,
	})
	if err != nil {
		t.Fatalf("CameraList search failed: %v", err)
	}
	if len(cameras) != 2 {
		t.Errorf("Expected 2 cameras matching search 'Door', got %d", len(cameras))
	}

	// Filter by type
	cameraType := models.CameraType("IP")
	ipCameras, err := cameraRepo.CameraList(&repository.CameraRepoFilter{
		Type: &cameraType,
	})
	if err != nil {
		t.Fatalf("CameraList filter type failed: %v", err)
	}
	if len(ipCameras) != 2 {
		t.Errorf("Expected 2 cameras of type IP, got %d", len(ipCameras))
	}

	// Pagination
	limit := 1
	offset := 1
	pagedCameras, err := cameraRepo.CameraList(&repository.CameraRepoFilter{
		Limit:  &limit,
		Offset: &offset,
	})
	if err != nil {
		t.Fatalf("CameraList pagination failed: %v", err)
	}
	if len(pagedCameras) != 1 {
		t.Errorf("Expected 1 camera in paginated list, got %d", len(pagedCameras))
	}
}
