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
data_info         = 'Data.Info'
device_response   = 'Device.Response'


## Subscribe topics
failure_network   = 'Failure.Network'
event_nac         = 'Event.NAC'
request_database  = 'Request.Database'
device_found      = 'Device.Found'
unauthorised_connection = 'Unauthorised.Connection'
device_request    = 'Device.Request'
device_update     = 'Device.Update'
status_nac        = 'Status.NAC'

## Bind
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=failure_network)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=event_nac)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=request_database)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=device_found)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=unauthorised_connection)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=device_request)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=device_update)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=status_nac)

print("Beginning Subscribe")
print("Waiting for notifications")
count = 1

def callback(ch, method, properties, body):
    print("Received %r:%r" % (method.routing_key, body))
    str = body.decode()
    print("NACIntegrator: I think we received a message: " + str)
    if method.routing_key == request_database:
        print("Request for data")
        request = { 
            "_id": count, 
            "_messageNum": 1,
            "_totalMessage": 1,
            "_topicMessage": "EVM3",
            "_timeSent": "12:00:30"
        }
        payload = json.dumps(request)
        channel.basic_publish(exchange='topics', routing_key=data_info, body=payload)
        print("Sent %r " % data_info)
        count = count + 1
    

channel.basic_consume(queue=queue_name, on_message_callback=callback, auto_ack=True)
channel.start_consuming()


