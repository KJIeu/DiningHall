import queue
import random
import threading
import time

import requests

from logic.Clients import orders, orders_bundle
from main import order_rating, orders_done
from model import tables_list, table_state

time_unit = 1

class Waiter(threading.Thread):
    def __init__(self, data, *args, **kwargs):
        super(Waiter, self).__init__(*args, **kwargs)
        self.id = data['id']
        self.name = data['name']
        self.daemon = True

    def run(self):
        while True:
            self.search_order()

    def search_order(self):
        try:
            order = orders.get()
            orders.task_done()
            table_id = next((i for i, table in enumerate(tables_list) if table['id'] == order['table_id']), None)
            print(
                f'{threading.current_thread().name} has taken the order with Id: {order["id"]} | priority: {order["priority"]} | items: {order["items"]} ')
            tables_list[table_id]['state'] = table_state[2]
            payload = dict({
                'order_id': order['id'],
                'table_id': order['table_id'],
                'waiter_id': self.id,
                'items': order['items'],
                'priority': order['priority'],
                'max_wait': order['max_wait'],
                'time_start': time.time()
            })

            time.sleep(random.randint(2, 4) * time_unit)
            # send order to kitchen
            requests.post('http://localhost:3030/order', json=payload, timeout=0.0000000001)

        except (queue.Empty, requests.exceptions.ReadTimeout, requests.exceptions.ConnectionError) as e:
            pass

    def serve_order(self, desired_order):
        # check if the order is the same order that what was requested
        received_order = next(
            (order for i, order in enumerate(orders_bundle) if order['id'] == desired_order['order_id']), None)
        if received_order is not None and received_order['items'].sort() == desired_order['items'].sort():
            # update table state
            table_id = next((i for i, table in enumerate(tables_list) if table['id'] == desired_order['table_id']),
                            None)
            tables_list[table_id]['state'] = table_state[3]
            order_serving_timestamp = int(time.time())
            order_pick_up_timestamp = int(desired_order['time_start'])
            # calculate total order time
            Order_total_preparing_time = order_serving_timestamp - order_pick_up_timestamp

            # calculate nr of start
            order_stars = {'order_id': desired_order['order_id']}
            if desired_order['max_wait'] > Order_total_preparing_time:
                order_stars['star'] = 5
            elif desired_order['max_wait'] * 1.1 > Order_total_preparing_time:
                order_stars['star'] = 4
            elif desired_order['max_wait'] * 1.2 > Order_total_preparing_time:
                order_stars['star'] = 3
            elif desired_order['max_wait'] * 1.3 > Order_total_preparing_time:
                order_stars['star'] = 2
            elif desired_order['max_wait'] * 1.4 > Order_total_preparing_time:
                order_stars['star'] = 1
            else:
                order_stars['star'] = 0

            order_rating.append(order_stars)
            sum_stars = sum(feedback['star'] for feedback in order_rating)
            avg = float(sum_stars / len(order_rating))

            served_order = {**desired_order, 'Serving_time': Order_total_preparing_time, 'status': 'DONE',
                            'Stars_feedback': order_stars}
            orders_done.append(served_order)
            # add to see the rating
            print(f'Serving the order Endpoint: /distribution :\n'
                  f'Order Id: {served_order["order_id"]}\n'
                  f'Waiter Id: {served_order["waiter_id"]}\n'
                  f'Table Id: {served_order["table_id"]}\n'
                  f'Items: {served_order["items"]}\n'
                  f'Priority: {served_order["priority"]}\n'
                  f'Max Wait: {served_order["max_wait"]}\n'
                  f'Waiting time: {served_order["Serving_time"]}\n'
                  f'Stars: {served_order["Stars_feedback"]}\n'
                  f'Restaurant rating: {avg}')

        else:
            raise Exception(
                f'Error. Provide the original order to costumer. Original: {received_order}, given: {desired_order}')
