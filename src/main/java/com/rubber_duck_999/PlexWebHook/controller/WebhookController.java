package com.rubber_duck_999.PlexWebHook.controller;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/plex")
public class WebhookController {

    @PostMapping
    public String handleWebhook(@RequestBody String payload) {
        // Log the payload received for debugging
        System.out.println("Received Webhook Payload: " + payload);

        // You can parse the payload and act on specific events
        // Handle your webhook logic here
        
        return "Webhook received successfully";
    }
}
