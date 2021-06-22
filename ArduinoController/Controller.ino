#include <ClickEncoder.h>
#include <TimerOne.h>

#include <LiquidCrystal.h>

#define LCD_RS       16
#define LCD_RW       9
#define LCD_EN      15
#define LCD_D4       5
#define LCD_D5       4
#define LCD_D6       3
#define LCD_D7       2

#define LCD_CHARS   16
#define LCD_LINES    2

#define LCD_DATA_LEN LCD_CHARS*LCD_LINES

LiquidCrystal lcd(LCD_RS, LCD_RW, LCD_EN, LCD_D4, LCD_D5, LCD_D6, LCD_D7);

ClickEncoder *encoder;
int16_t last, value;

long last_send = 0;


void timerIsr() {
  encoder->service();
}



void setup() {
  Serial.begin(9600);
  encoder = new ClickEncoder(A1, A0, A2);
  encoder->setAccelerationEnabled(false);

  lcd.begin(LCD_CHARS, LCD_LINES);
  lcd.clear();

  Timer1.initialize(1000);
  Timer1.attachInterrupt(timerIsr);

  last = -1;
}

char data[LCD_DATA_LEN] = "0123456789abcdef0123456789abcdef";
char data1[LCD_CHARS];
char data2[LCD_CHARS];


void printCelcius() {
  int val = analogRead(A3);
  float mv = ( val / 1024.0) * 5000;
  float cel = mv / 10;
  Serial.print(cel);
}

void loop() {
  if (Serial.available() >= LCD_DATA_LEN) {
    int n = Serial.readBytes(data, LCD_DATA_LEN);
    //Clear Buffer
    while (Serial.available() > 0) {
      Serial.read();
    }
    //Serial.println(n);
    //copy halfs
    for (int i = 0; i < LCD_CHARS; i++) {
      data1[i] = data[i];
      data2[i] = data[i + LCD_CHARS];
    }
    lcd.setCursor(0, 0);
    lcd.println(data1);
    lcd.setCursor(0, 1);
    lcd.println(data2);

  }



  value += encoder->getValue();
  if (value != last) {
    if (value > last) {
      Serial.println("go-up");
    } else if (value < last) {
      Serial.println("go-down");
    }
    last_send = millis();
    last = value;
  }
  ClickEncoder::Button b = encoder->getButton();

  if (b != ClickEncoder::Open && b == ClickEncoder::Clicked) {
    Serial.println("btn-clk");
    last_send = millis();
  }

  if (abs(millis() - last_send) > 1000) {
    Serial.print("update ");
    printCelcius();
    Serial.println();
    last_send = millis();
  }
}
/*
  if (b != ClickEncoder::Open) {
  Serial.print("Button: ");
  #define VERBOSECASE(label) case label: Serial.println(#label); break;
  switch (b) {
      VERBOSECASE(ClickEncoder::Pressed);
      VERBOSECASE(ClickEncoder::Held)
      VERBOSECASE(ClickEncoder::Released)
      VERBOSECASE(ClickEncoder::Clicked)
    case ClickEncoder::DoubleClicked:
      Serial.println("DoubleClicked");

      break;
  }
  }
*/
