'''
Created on 30 Dec 2019

@author: Rubber-Duck-999
'''

#!/usr/bin/env python
import pika, sys, time, json


print("## NAC Integrator Start up")
credentials = pika.PlainCredentials('guest', 'password')
connection = pika.BlockingConnection(pika.ConnectionParameters('localhost', 5672, '/', credentials))
channel = connection.channel()
channel.exchange_declare(exchange='topics', exchange_type='topic', durable=True)
#
result = channel.queue_declare(queue='', exclusive=False, durable=True)
queue_name = result.method.queue
#

## Publish Topics
device_response   = 'Device.Response'


## Subscribe topics
failure_network   = 'Failure.Network'
device_found      = 'Device.Found'
device_request    = 'Device.Request'
status_request    = 'Status.Request'

## Bind
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=failure_network)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=device_found)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=device_request)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=status_request)

print("Beginning Subscribe")
print("Waiting for notifications")
count = 1

def callback(ch, method, properties, body):
    str = body.decode()
    message = json.loads(str)
    print(str)
    if method.routing_key == device_request:
        request = { 
            "id": message['id'],
            "name": message['name'],
            "mac": message['mac'],
            "status": 'BLOCKED'
        }
        payload = json.dumps(request)
        channel.basic_publish(exchange='topics', routing_key=device_response, body=payload)
    

channel.basic_consume(queue=queue_name, on_message_callback=callback, auto_ack=True)
channel.start_consuming()


