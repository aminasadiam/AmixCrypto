import pandas as pd
import sys
from sklearn.ensemble import RandomForestRegressor
from sklearn.model_selection import train_test_split

def main():
    interval = sys.argv[2]
    if interval == "1h":
        Train1H()
    elif interval == "4h":
        Train4H()
    elif interval == "1d":
        Train1D()
    else:
        print("Invalid Interval!")


def Train1H():
    symbol = sys.argv[1]
    path = f"./Data/{symbol}.csv"

    data = pd.read_csv(path)
    X = data[["Open","High","Low","Volume","RSI","SMA"]]
    y = data['Close']

    X_train, X_test, y_train, y_test = train_test_split(X, y, train_size=0.8, test_size=0.2, random_state=42)

    model = RandomForestRegressor()
    model.fit(X_train, y_train)

    last_data = data.iloc[-1:]  # آخرین داده‌ها برای پیش‌بینی
    future_predictions = []
    times = pd.to_datetime(data['Time'])
    future_dates = pd.date_range(start=times.max() + pd.Timedelta(hours=1), periods=5, freq='h')

    for i in range(5):
        prediction = model.predict(last_data[["Open","High","Low","Volume","RSI","SMA"]])
        future_predictions.append(prediction[i])

        # به‌روز‌رسانی داده‌های ورودی برای پیش‌بینی بعدی
        new_row = {
            'Open': prediction[i],
            'High': prediction[i],
            'Low': prediction[i],
            'Volume': last_data['Volume'].values[i],
            'SMA': last_data['SMA'].values[i],
            'RSI': last_data['RSI'].values[i],
            'Time': pd.to_datetime(last_data['Time'].values[i]) + pd.Timedelta(hours=1)
        }

        last_data = pd.concat([last_data, pd.DataFrame([new_row])], ignore_index=True)

    actuals = [data[data['Time'] == future_dates[i] - pd.Timedelta(hours=1)]['Close'].values[0] if not data[data['Time'] == future_dates[i] - pd.Timedelta(hours=1)].empty else None for i in range(5)]
    result = {
        "accuracy": "N/A",  # دقت مدل برای پیش‌بینی‌های آینده را نمی‌توان به سادگی محاسبه کرد
        "predictions": future_predictions,
        "dates": future_dates.strftime('%m/%d/%Y, %H:%M:%S').tolist(),  # تبدیل به فرمت تاریخ مناسب
        "actuals": actuals
    }

    for i in range(5):
        print(f"{result['dates'][i]} -> {result['predictions'][i]}")

def Train4H():
    symbol = sys.argv[1]
    path = f"./Data/{symbol}.csv"

    data = pd.read_csv(path)
    X = data[["Open","High","Low","Volume","RSI","SMA"]]
    y = data['Close']

    X_train, X_test, y_train, y_test = train_test_split(X, y, train_size=0.8, test_size=0.2, random_state=42)

    model = RandomForestRegressor()
    model.fit(X_train, y_train)

    last_data = data.iloc[-1:]  # آخرین داده‌ها برای پیش‌بینی
    future_predictions = []
    times = pd.to_datetime(data['Time'])
    future_dates = pd.date_range(start=times.max() + pd.Timedelta(hours=4), periods=5, freq='4h')

    for i in range(5):
        prediction = model.predict(last_data[["Open","High","Low","Volume","RSI","SMA"]])
        future_predictions.append(prediction[i])

        # به‌روز‌رسانی داده‌های ورودی برای پیش‌بینی بعدی
        new_row = {
            'Open': prediction[i],
            'High': prediction[i],
            'Low': prediction[i],
            'Volume': last_data['Volume'].values[i],
            'SMA': last_data['SMA'].values[i],
            'RSI': last_data['RSI'].values[i],
            'Time': pd.to_datetime(last_data['Time'].values[i]) + pd.Timedelta(hours=4)
        }

        last_data = pd.concat([last_data, pd.DataFrame([new_row])], ignore_index=True)

    actuals = [data[data['Time'] == future_dates[i] - pd.Timedelta(hours=4)]['Close'].values[0] if not data[data['Time'] == future_dates[i] - pd.Timedelta(hours=4)].empty else None for i in range(5)]
    result = {
        "accuracy": "N/A",  # دقت مدل برای پیش‌بینی‌های آینده را نمی‌توان به سادگی محاسبه کرد
        "predictions": future_predictions,
        "dates": future_dates.strftime('%m/%d/%Y, %H:%M:%S').tolist(),  # تبدیل به فرمت تاریخ مناسب
        "actuals": actuals
    }

    for i in range(5):
        print(f"{result['dates'][i]} -> {result['predictions'][i]}")

def Train1D():
    symbol = sys.argv[1]
    path = f"./Data/{symbol}.csv"

    data = pd.read_csv(path)
    X = data[["Open","High","Low","Volume","RSI","SMA"]]
    y = data['Close']

    X_train, X_test, y_train, y_test = train_test_split(X, y, train_size=0.8, test_size=0.2, random_state=42)

    model = RandomForestRegressor()
    model.fit(X_train, y_train)

    last_data = data.iloc[-1:]  # آخرین داده‌ها برای پیش‌بینی
    future_predictions = []
    times = pd.to_datetime(data['Time'])
    future_dates = pd.date_range(start=times.max() + pd.Timedelta(days=1), periods=5, freq='D')

    for i in range(5):
        prediction = model.predict(last_data[["Open","High","Low","Volume","RSI","SMA"]])
        future_predictions.append(prediction[i])

        # به‌روز‌رسانی داده‌های ورودی برای پیش‌بینی بعدی
        new_row = {
            'Open': prediction[i],
            'High': prediction[i],
            'Low': prediction[i],
            'Volume': last_data['Volume'].values[i],
            'SMA': last_data['SMA'].values[i],
            'RSI': last_data['RSI'].values[i],
            'Time': pd.to_datetime(last_data['Time'].values[i]) + pd.Timedelta(days=1)
        }

        last_data = pd.concat([last_data, pd.DataFrame([new_row])], ignore_index=True)

    actuals = [data[data['Time'] == future_dates[i] - pd.Timedelta(days=1)]['Close'].values[0] if not data[data['Time'] == future_dates[i] - pd.Timedelta(days=1)].empty else None for i in range(5)]
    result = {
        "accuracy": "N/A",  # دقت مدل برای پیش‌بینی‌های آینده را نمی‌توان به سادگی محاسبه کرد
        "predictions": future_predictions,
        "dates": future_dates.strftime('%m/%d/%Y, %H:%M:%S').tolist(),  # تبدیل به فرمت تاریخ مناسب
        "actuals": actuals
    }

    for i in range(5):
        print(f"{result['dates'][i]} -> {result['predictions'][i]}")

if __name__ == "__main__":
    main()