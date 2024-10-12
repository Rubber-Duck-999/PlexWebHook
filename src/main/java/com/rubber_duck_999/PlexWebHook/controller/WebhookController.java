package com.rubber_duck_999.PlexWebHook.controller;

import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.http.ResponseEntity;
import org.springframework.http.HttpStatus;

import java.util.Map;

@RestController
@RequestMapping("/webhook")
public class WebhookController {

    @PostMapping(consumes = "multipart/form-data")
    public ResponseEntity<String> handleMultipartWebhook(
            @RequestParam Map<String, String> formData,
            @RequestPart("file") MultipartFile file) {

        // Log the received form data (payload)
        System.out.println("Received Form Data: " + formData);

        // Extract specific form fields from the payload
        String eventType = formData.get("eventType");
        String timestamp = formData.get("timestamp");

        System.out.println("Event Type: " + eventType);
        System.out.println("Timestamp: " + timestamp);

        // Handle file upload if present
        if (file != null && !file.isEmpty()) {
            try {
                String fileContent = new String(file.getBytes());
                System.out.println("File content: " + fileContent);
            } catch (Exception e) {
                return new ResponseEntity<>("Failed to process file", HttpStatus.INTERNAL_SERVER_ERROR);
            }
        }

        return new ResponseEntity<>("Webhook received successfully", HttpStatus.OK);
    }
}
