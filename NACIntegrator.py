'''
Created on 30 Dec 2019

@author: Rubber-Duck-999
'''

#!/usr/bin/env python
import pika, sys, time


print("## NAC Integrator Start up")
connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
channel = connection.channel()
#
channel.exchange_declare(exchange='topics', exchange_type='topic', durable=True)
#
result = channel.queue_declare(queue='', exclusive=False, durable=True)
queue_name = result.method.queue
#

## Topics
request_image     = 'Request.Image'
request_data      = 'Request.Data'
request_access    = 'Request.Access'
request_weather   = 'Request.Weather'
data_info         = 'Data.Info'
failure_network   = 'Failure.Network'
event_nac         = 'Event.NAC'
request_database  = 'Request.Database'
data_response     = 'Data.Response'
access_response   = 'Access.Response'
failure_component = 'Failure.Component'
weather           = 'Weather'

channel.queue_bind(exchange='topics', queue=queue_name, routing_key=failure_network)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=failure_component)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=event_nac)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=request_database)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=data_response)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=access_response)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=weather)

print("Beginning Subscribe")
print("Waiting for notifications")
count = 0

def callback(ch, method, properties, body):
    print("Received: " % (method.routing_key, body))
    str = body.decode()
    print("NACIntegrator: I think we received a message: " + str)
    count = count + 1
    print("Publishing " + count)
    text = text_to_send + ' ' + str(count)
    channel.basic_publish(exchange='topics', routing_key=request_image, body=text)
    time.sleep(0.5)
    channel.basic_publish(exchange='topics', routing_key=request_data, body=text)
    time.sleep(0.5)
    channel.basic_publish(exchange='topics', routing_key=request_access, body=text)
    time.sleep(0.5)
    channel.basic_publish(exchange='topics', routing_key=request_weather, body=text)
    time.sleep(0.5)
    channel.basic_publish(exchange='topics', routing_key=data_info, body=text)
    

channel.basic_consume(queue=queue_name, on_message_callback=callback, auto_ack=True)
channel.start_consuming()


