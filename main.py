from flask import Flask, request
import threading

from logic import Waiter
from logic.Clients import Clients

from model import waiters
from model import tables_list
from model import table_state

orders_done = []
order_rating = []

app = Flask(__name__)
threads = []




@app.route('/distribution', methods=['POST'])
def distribution():
    order = request.get_json()
    print(f'Received order from kitchen. Order ID: {order["order_id"]} items: {order["items"]} ')
    table_id = next((i for i, table in enumerate(tables_list) if table['id'] == order['table_id']), None)
    tables_list[table_id]['state'] = table_state[2]
    # get the waiter thread and serve order
    waiter_thread: Waiter = next((w for w in threads if type(w) == Waiter and w.id == order['waiter_id']), None)
    waiter_thread.serve_order(order)  #
    return {'isSuccess': True}




def run_dinninghall_server():
    main_thread = threading.Thread(target=lambda: app.run(host='0.0.0.0', port=3000, debug=False, use_reloader=False),
                                   daemon=True)
    threads.append(main_thread)
    client_thread = Clients()
    threads.append(client_thread)
    for _, w in enumerate(waiters):
        waiter_thread = Waiter(w)
        threads.append(waiter_thread)
    for th in threads:
        th.start()
    for th in threads:
        th.join()


if __name__ == '__main__':
    run_dinninghall_server()