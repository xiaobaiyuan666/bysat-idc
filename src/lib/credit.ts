import { type CreditTransactionType, Prisma } from "@prisma/client";

import { db } from "@/lib/db";

type DbClient = Prisma.TransactionClient;

export async function changeCustomerBalance(input: {
  customerId: string;
  amount: number;
  type: CreditTransactionType;
  description: string;
  operatorId?: string | null;
  allowNegative?: boolean;
  tx?: DbClient;
}) {
  const client = input.tx ?? db;

  const customer = await client.customer.findUnique({
    where: {
      id: input.customerId,
    },
  });

  if (!customer) {
    throw new Error("CUSTOMER_NOT_FOUND");
  }

  const nextBalance = customer.creditBalance + input.amount;

  if (nextBalance < 0 && !input.allowNegative) {
    throw new Error("INSUFFICIENT_BALANCE");
  }

  const updatedCustomer = await client.customer.update({
    where: {
      id: customer.id,
    },
    data: {
      creditBalance: nextBalance,
    },
  });

  const transaction = await client.creditTransaction.create({
    data: {
      customerId: customer.id,
      operatorId: input.operatorId ?? undefined,
      type: input.type,
      amount: input.amount,
      balanceAfter: nextBalance,
      description: input.description,
    },
  });

  return {
    customer: updatedCustomer,
    transaction,
  };
}
